WEB_CLIENT_NAME=sales-tracker_web-client

.PHONY: web-client-create web-client-deps-init

web-client-create:
	pnpm create next-app@latest $(WEB_CLIENT_NAME) --typescript --eslint --react-compiler --no-tailwind --no-src-dir --app --import-alias "@/*" --use-pnpm

web-client-deps-init:
	cd $(WEB_CLIENT_NAME) && pnpm add @emotion/cache @emotion/react @emotion/styled @emotion/sheet @emotion/weak-memoize @emotion/memoize @emotion/serialize @emotion/unitless @emotion/hash @mui/icons-material @mui/material @mui/material-nextjs @mui/system @mui/private-theming @mui/styled-engine react react-dom next react-hook-form && pnpm add -D typescript @types/node @types/react @types/react-dom babel-plugin-react-compiler eslint eslint-config-next stylis
