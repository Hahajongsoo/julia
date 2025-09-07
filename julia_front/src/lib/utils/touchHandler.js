// 터치/스와이프 이벤트 처리 유틸리티 모듈

// 터치 상태 관리 클래스
export class TouchState {
	constructor() {
		this.touchStartX = 0;
		this.touchStartY = 0;
		this.touchEndX = 0;
		this.touchEndY = 0;
		this.isSwiping = false;
	}

	reset() {
		this.touchStartX = 0;
		this.touchStartY = 0;
		this.touchEndX = 0;
		this.touchEndY = 0;
		this.isSwiping = false;
	}
}

// 터치 시작 이벤트 핸들러
export function handleTouchStart(event, touchState) {
	// Svelte 커스텀 이벤트인 경우 event.detail에서 실제 터치 이벤트 추출
	const touchEvent = event.detail || event;
	
	// 터치 이벤트가 아니거나 터치 포인트가 없는 경우 방어
	if (!touchEvent.touches || touchEvent.touches.length === 0) {
		return;
	}
	
	touchState.touchStartX = touchEvent.touches[0].clientX;
	touchState.touchStartY = touchEvent.touches[0].clientY;
	touchState.isSwiping = false;
}

// 터치 이동 이벤트 핸들러
export function handleTouchMove(event, touchState) {
	if (!touchState.touchStartX || !touchState.touchStartY) return;
	
	// Svelte 커스텀 이벤트인 경우 event.detail에서 실제 터치 이벤트 추출
	const touchEvent = event.detail || event;
	
	// 터치 이벤트가 아니거나 터치 포인트가 없는 경우 방어
	if (!touchEvent.touches || touchEvent.touches.length === 0) {
		return;
	}

	touchState.touchEndX = touchEvent.touches[0].clientX;
	touchState.touchEndY = touchEvent.touches[0].clientY;

	const deltaX = touchState.touchStartX - touchState.touchEndX;
	const deltaY = touchState.touchStartY - touchState.touchEndY;

	// 수평 스와이프가 수직 스와이프보다 클 때만 스와이프로 인식
	if (Math.abs(deltaX) > Math.abs(deltaY) && Math.abs(deltaX) > 50) {
		touchState.isSwiping = true;
	}
}

// 터치 종료 이벤트 핸들러
export function handleTouchEnd(event, touchState, callbacks = {}) {
	// Svelte 커스텀 이벤트인 경우 event.detail에서 실제 터치 이벤트 추출
	const touchEvent = event.detail || event;
	
	// 터치 이벤트가 아니거나 터치 포인트가 없는 경우 방어
	if (!touchEvent.changedTouches || touchEvent.changedTouches.length === 0) {
		touchState.reset();
		return;
	}
	
	if (!touchState.isSwiping || !touchState.touchStartX || !touchState.touchEndX) {
		touchState.reset();
		return;
	}

	const deltaX = touchState.touchStartX - touchState.touchEndX;
	const minSwipeDistance = 50;

	if (Math.abs(deltaX) > minSwipeDistance) {
		if (deltaX > 0) {
			// 왼쪽으로 스와이프 - 다음 달
			if (callbacks.onNextMonth) {
				callbacks.onNextMonth();
			}
		} else {
			// 오른쪽으로 스와이프 - 이전 달
			if (callbacks.onPreviousMonth) {
				callbacks.onPreviousMonth();
			}
		}
	}

	// 터치 상태 초기화
	touchState.reset();
}

// 터치 이벤트 핸들러들을 생성하는 팩토리 함수
export function createTouchHandlers(touchState, callbacks = {}) {
	return {
		onTouchStart: (event) => handleTouchStart(event, touchState),
		onTouchMove: (event) => handleTouchMove(event, touchState),
		onTouchEnd: (event) => handleTouchEnd(event, touchState, callbacks)
	};
}
