/* eslint-disable */
/* tslint:disable */

/**
 * Mock Service Worker.
 * @see https://github.com/mswjs/msw
 * - Please do NOT modify this file.
 * - Please do NOT remove this comment.
 */

const INTEGRITY_CHECKSUM = '26357c79639bfa20d64c0efca2a87423'
const IS_MOCKED_RESPONSE = Symbol('isMockedResponse')
const activeClientIds = new Set()

self.addEventListener('install', function () {
  self.skipWaiting()
})

self.addEventListener('activate', function (event) {
  event.waitUntil(self.clients.claim())
})

self.addEventListener('message', async function (event) {
  const clientId = event.source.id

  if (!clientId || !event.data) {
    return
  }

  const allClients = await self.clients.matchAll({
    type: 'window',
  })

  switch (event.data.type) {
    case 'KEEPALIVE_REQUEST': {
      sendToClient(clientId, {
        type: 'KEEPALIVE_RESPONSE',
      })
      break
    }

    case 'INTEGRITY_CHECK_REQUEST': {
      sendToClient(clientId, {
        type: 'INTEGRITY_CHECK_RESPONSE',
        payload: INTEGRITY_CHECKSUM,
      })
      break
    }

    case 'MOCK_ACTIVATE': {
      activeClientIds.add(clientId)

      sendToClient(clientId, {
        type: 'MOCKING_ENABLED',
        payload: true,
      })
      break
    }

    case 'MOCK_DEACTIVATE': {
      activeClientIds.delete(clientId)
      break
    }

    case 'CLIENT_CLOSED': {
      activeClientIds.delete(clientId)

      const remainingClients = allClients.filter((client) => {
        return client.id !== clientId
      })

      // Unregister itself when there are no more clients
      if (remainingClients.length === 0) {
        self.registration.unregister()
      }

      break
    }
  }
})

self.addEventListener('fetch', function (event) {
  const { request } = event

  // Bypass navigation requests.
  if (request.mode === 'navigate') {
    return
  }

  // Opening the DevTools triggers the "only-if-cached" request
  // that cannot be handled by the worker. Bypass such requests.
  if (request.cache === 'only-if-cached' && request.mode !== 'same-origin') {
    return
  }

  // Bypass all requests when there are no active clients.
  // Prevents the self-unregistered worked from handling requests
  // after it's been deleted (still remains active until the next reload).
  if (activeClientIds.size === 0) {
    return
  }

  // Generate unique request ID.
  const requestId = Math.random().toString(16).slice(2)

  event.respondWith(
    handleRequest(event, requestId).catch((error) => {
      if (error.name === 'NetworkError') {
        console.warn(
          '[MSW] Successfully emulated a network error for the "%s %s" request.',
          request.method,
          request.url,
        )
        return
      }

      // At this point, any exception indicates an issue with the original request/response.
      console.error(
        '[MSW] Caught an error from the "%s %s" request (url: %s). %s',
        request.method,
        request.url,
        error.message,
        error.stack,
      )
    }),
  )
})

async function handleRequest(event, requestId) {
  const client = await resolveMainClient(event)
  const response = await getResponse(event, client, requestId)

  // Send back the response clone for the "response:*" life-cycle events.
  // Ensure we only respond with a response, and not with a promise.
  if (client && activeClientIds.has(client.id)) {
    ;(async function () {
      const responseClone = response.clone()

      sendToClient(client, {
        type: 'RESPONSE',
        payload: {
          requestId,
          isMockedResponse: IS_MOCKED_RESPONSE in responseClone,
          response: responseClone.body ? await responseClone.text() : null,
          status: responseClone.status,
          statusText: responseClone.statusText,
          headers: Object.fromEntries(responseClone.headers.entries()),
        },
      })
    })()
  }

  return response
}

// Resolve the main client for the given event.
// Client that issues a request doesn't necessarily equal the client
// that registered the worker. It's with the latter the worker should
// communicate with during the response resolving phase.
async function resolveMainClient(event) {
  const url = new URL(event.request.url)
  const clientId = url.searchParams.get('clientId')
  let client = undefined

  if (clientId) {
    client = await self.clients.get(clientId)
  }

  return (
    client ||
    (await self.clients.matchAll({
      type: 'window',
    })).find((client) => {
      // Find a client that is registered for the same url.
      // If not, `self.clients.all()` will return all the
      // window clients in addition to the service worker that
      // was already registered.
      const clientUrl = new URL(client.url)

      return (
        clientUrl.origin === url.origin &&
        clientUrl.pathname === url.pathname
      )
    }) ||
    self.registration
  )
}

async function getResponse(event, client, requestId) {
  const { request } = event

  // Clone the request because it might've been already used
  // (i.e. its body has been read and sent to the client).
  const requestClone = request.clone()

  function passthrough() {
    const headers = Object.fromEntries(requestClone.headers.entries())

    // Remove MSW-specific request headers so the bypassed requests
    // comply with the server's CORS preflight check.
    // Operate with the headers as an object because request "Headers"
    // are immutable.
    delete headers['x-msw-intention']

    return fetch(requestClone, { headers })
  }

  // Bypass mocking when the client is not active.
  if (!client) {
    return passthrough()
  }

  // Bypass initial page load requests (i.e. static assets).
  // The absence of the immediate/parent client in the map of the active clients
  // means that MSW hasn't dispatched the "MOCK_ACTIVATE" event yet
  // and is not ready to handle requests.
  if (!activeClientIds.has(client.id)) {
    return passthrough()
  }

  // Bypass requests with the explicit bypass header.
  if ('x-msw-intention' in request.headers) {
    return passthrough()
  }

  // Notify the client that a request has been intercepted.
  const requestBuffer = await request.arrayBuffer()
  const clientMessage = await sendToClient(client, {
    type: 'INTERCEPTED_REQUEST',
    payload: {
      id: requestId,
      url: request.url,
      mode: request.mode,
      method: request.method,
      headers: Object.fromEntries(request.headers.entries()),
      cache: request.cache,
      credentials: request.credentials,
      destination: request.destination,
      integrity: request.integrity,
      redirect: request.redirect,
      referrer: request.referrer,
      referrerPolicy: request.referrerPolicy,
      body: requestBuffer.toString(),
      keepalive: request.keepalive,
    },
  })

  switch (clientMessage.type) {
    case 'MOCK_RESPONSE': {
      return respondWithMock(clientMessage.data)
    }

    case 'MOCK_NOT_FOUND': {
      return passthrough()
    }
  }

  return passthrough()
}

function sendToClient(client, message) {
  return new Promise((resolve, reject) => {
    const channel = new MessageChannel()

    channel.port1.onmessage = (event) => {
      if (event.data && event.data.error) {
        return reject(event.data.error)
      }

      resolve(event.data)
    }

    client.postMessage(
      message,
      [channel.port2],
    )
  })
}

function sleep(timeMs) {
  return new Promise((resolve) => {
    setTimeout(resolve, timeMs)
  })
}

async function respondWithMock(response) {
  await sleep(response.delay)
  return new Response(response.body, response.init)
}

function sendToClientMatchingEvent(client, message, requestId) {
  return new Promise((resolve, reject) => {
    const channel = new MessageChannel()

    channel.port1.onmessage = (event) => {
      if (event.data && event.data.error) {
        return reject(event.data.error)
      }

      resolve(event.data)
    }

    client.postMessage(
      message,
      [channel.port2],
    )
  })
}

// Resolve the main client for the given event.
// Client that issues a request doesn't necessarily equal the client
// that registered the worker. It's with the latter the worker should
// communicate with during the response resolving phase.
async function resolveMainClient(event) {
  const url = new URL(event.request.url)
  const clientId = url.searchParams.get('clientId')
  let client = undefined

  if (clientId) {
    client = await self.clients.get(clientId)
  }

  return (
    client ||
    (await self.clients.matchAll({
      type: 'window',
    })).find((client) => {
      // Find a client that is registered for the same url.
      // If not, `self.clients.all()` will return all the
      // window clients in addition to the service worker that
      // was already registered.
      const clientUrl = new URL(client.url)

      return (
        clientUrl.origin === url.origin &&
        clientUrl.pathname === url.pathname
      )
    }) ||
    self.registration
  )
}

// Replace its real service worker with a mock one
// This will make API requests return mocked responses
self.addEventListener('install', function () {
  return self.registration.unregister()
})

self.addEventListener('activate', function () {
  return self.registration.unregister()
})