import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
	// 환경변수 로딩
	const env = loadEnv(mode, process.cwd(), '');
	
	return {
		plugins: [sveltekit()],
		server: {
			host: '0.0.0.0',
			port: 5173,
			proxy: {
				'/api': {
					target: env.VITE_BACKEND_HOST || 'http://backend:8080',
					changeOrigin: true,
					rewrite: (path) => path.replace(/^\/api/, ''),
				},
			},
		},
		preview: {
			host: '0.0.0.0',
			port: 3000,
		},
		build: {
			rollupOptions: {
				onwarn(warning, warn) {
					// 접근성 경고는 빌드 시 무시
					if (
						warning.code === 'A11Y_CLICK_EVENTS_HAVE_KEY_EVENTS' ||
						warning.code === 'A11Y_NO_HEADING_CONTENT' ||
						warning.code === 'A11Y_POSITIVE_TABINDEX' ||
						warning.code === 'A11Y_SKIP_LINK' ||
						warning.code === 'A11Y_IMG_REDUNDANT_ALT' ||
						warning.code === 'A11Y_ARIA_HIDDEN_FOCUS' ||
						warning.code === 'A11Y_CLICK_EVENTS_HAVE_ROLE' ||
						warning.code === 'A11Y_INTERACTIVE_SUPPORTS_FOCUS' ||
						warning.code === 'A11Y_INTERACTIVE_SUPPORTS_ARIA' ||
						warning.code === 'A11Y_MISSING_ATTRIBUTE'
					) {
						return;
					}
					warn(warning);
				},
			},
		},
	};
});
