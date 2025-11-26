WEB_CLIENT_NAME=warehouse-control_web-client

.PHONY: web-client-create web-client-deps-init

web-client-create:
	pnpm create next-app@latest $(WEB_CLIENT_NAME) \
	--typescript \
	--eslint \
	--react-compiler \
	--no-tailwind \
	--src-dir \
	--app \
	--turbopack \
	--import-alias "@/*" \
	--use-pnpm

web-client-deps-init:
	cd $(WEB_CLIENT_NAME) && \
	pnpm add @mui/material @emotion/react @emotion/styled && \
	pnpm add @fontsource/roboto && \
	pnpm add @mui/icons-material && \
	pnpm add @mui/material-nextjs \
	pnpm add material-react-table @mui/x-date-pickers