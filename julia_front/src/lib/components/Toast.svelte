<script>
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	// 내부 상태로 관리
	let message = '';
	let type = 'info'; // 'info', 'success', 'error'
	let duration = 3000;

	let toastElement;
	let autoHideTimeout;

	// 토스트 표시
	function showToast(msg, toastType = 'info', toastDuration = 3000) {
		// 이전 자동 제거 타이머가 있으면 취소
		if (autoHideTimeout) {
			clearTimeout(autoHideTimeout);
		}
		
		message = msg;
		type = toastType;
		duration = toastDuration;
		
		// 토스트 요소가 DOM에 추가되면 애니메이션 시작
		setTimeout(() => {
			if (toastElement) {
				toastElement.style.animation = 'slideInRight 0.3s ease-out';
			}
		}, 10);

		// 자동 제거
		autoHideTimeout = setTimeout(() => {
			hideToast();
		}, duration);
	}

	// 토스트 숨기기
	function hideToast() {
		// 자동 제거 타이머 취소
		if (autoHideTimeout) {
			clearTimeout(autoHideTimeout);
			autoHideTimeout = null;
		}
		
		if (toastElement) {
			toastElement.style.animation = 'slideOutRight 0.3s ease-in forwards';
			setTimeout(() => {
				message = ''; // 메시지 초기화
				dispatch('close');
			}, 300);
		}
	}

	// 외부에서 호출할 수 있도록 함수 노출
	export { showToast, hideToast };
</script>

{#if message}
	<div 
		bind:this={toastElement}
		class="toast toast-{type}"
	>
		{message}
	</div>
{/if}

<style>
	.toast {
		position: fixed;
		top: 20px;
		right: 20px;
		padding: 12px 20px;
		border-radius: 8px;
		z-index: 10000;
		font-size: 14px;
		font-weight: 500;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
		color: white;
		animation: slideInRight 0.3s ease-out;
	}

	.toast-success {
		background: #10b981;
	}

	.toast-error {
		background: #ef4444;
	}

	.toast-info {
		background: #3b82f6;
	}

	@keyframes slideInRight {
		from {
			transform: translateX(100%);
			opacity: 0;
		}
		to {
			transform: translateX(0);
			opacity: 1;
		}
	}

	@keyframes slideOutRight {
		from {
			transform: translateX(0);
			opacity: 1;
		}
		to {
			transform: translateX(100%);
			opacity: 0;
		}
	}
</style>
