WEB_CLIENT_NAME=image-processor_web-client

.PHONY: web-client-create web-client-deps-init

web-client-create:
	pnpm create next-app@latest $(WEB_CLIENT_NAME) --typescript --eslint --react-compiler --no-tailwind --no-src-dir --app --import-alias "@/*" --use-pnpm

web-client-deps-init:
	cd $(WEB_CLIENT_NAME) && pnpm add @emotion/cache @emotion/react @emotion/styled @mui/icons-material @mui/material @mui/material-nextjs react react-dom next react-hook-form && pnpm add -D typescript @types/node @types/react @types/react-dom babel-plugin-react-compiler eslint eslint-config-next
