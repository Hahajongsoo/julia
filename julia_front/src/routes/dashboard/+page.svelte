<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { isAuthenticated, logout, fetchWithAuth } from '$lib/auth';
	import { API_ENDPOINTS } from '$lib/config';

	// SvelteKit props (현재 사용하지 않음)
	// export let data;
	// export let params;

	let username = '';
	let userRole = 'student';
	let isLoading = false;
	let userLoading = false;
	let students = [];
	let studentsLoading = false;
	let classes = [];
	let classesLoading = false;
	let expandedClasses = new Set(); // 토글된 반들을 추적

	// 시험 기간 관련 상태
	let examPeriods = [];
	let examPeriodsLoading = false;
	let selectedClassForExam = null;
	let showExamPeriods = false;
	let allExamPeriods = []; // 전체 시험 기간 데이터
	let examPeriodsSummary = {
		total: 0,
		upcoming: 0,
		ongoing: 0,
		completed: 0
	};

	// 학생 생성/수정 관련 상태
	let showCreateForm = false;
	let showEditForm = false;
	let selectedStudent = null;
	let isCreating = false;
	let isEditing = false;

	// 반 생성/수정 관련 상태
	let showCreateClassForm = false;
	let showEditClassForm = false;
	let selectedClass = null;
	let isCreatingClass = false;
	let isEditingClass = false;

	// 시험 기간 생성/수정 관련 상태
	let showCreateExamPeriodForm = false;
	let showEditExamPeriodForm = false;
	let selectedExamPeriod = null;
	let isCreatingExamPeriod = false;
	let isEditingExamPeriod = false;

	let createFormData = {
		id: '',
		password: '',
		phone: '',
		class_id: null,
		role: 'student',
	};

	let editFormData = {
		id: '',
		password: '',
		phone: '',
		class_id: null,
		role: 'student',
	};

	let createClassFormData = {
		class_name: '',
	};

	let editClassFormData = {
		class_id: null,
		class_name: '',
	};

	let createExamPeriodFormData = {
		class_id: null,
		name: '',
		description: '',
		start_date: '',
		end_date: '',
		english_date: '',
	};

	let editExamPeriodFormData = {
		exam_period_id: null,
		class_id: null,
		name: '',
		description: '',
		start_date: '',
		end_date: '',
		english_date: '',
	};

	onMount(async () => {
		console.log('대시보드 페이지 로드됨');
		
		// 로그인 상태 확인
		const authenticated = await isAuthenticated();
		console.log('인증 상태 확인:', authenticated);
		
		if (!authenticated) {
			console.log('인증되지 않음, 로그인 페이지로 이동');
			await goto('/');
			return;
		}

		console.log('인증됨, 사용자 정보 로드 중...');
		// 현재 사용자 정보 가져오기
		await loadCurrentUser();

		// 관리자인 경우에만 학생 목록과 반 목록, 시험 기간 요약 로드
		if (userRole === 'admin') {
			await Promise.all([loadStudents(), loadClasses(), loadAllExamPeriods()]);
		}
	});

	async function loadCurrentUser() {
		userLoading = true;
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.ME);
			if (response.ok) {
				const userData = await response.json();
				username = userData.id || '사용자';
				userRole = userData.role || 'student';
			} else {
				console.error('Failed to load user data');
				username = '사용자';
				userRole = 'student';
			}
		} catch (error) {
			console.error('Error loading user data:', error);
			username = '사용자';
			userRole = 'student';
		} finally {
			userLoading = false;
		}
	}

	async function loadStudents() {
		studentsLoading = true;
		try {
			const response = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/users`);
			if (response.ok) {
				students = await response.json();
			} else {
				console.error('Failed to load students');
				students = [];
			}
		} catch (error) {
			console.error('Error loading students:', error);
			students = [];
		} finally {
			studentsLoading = false;
		}
	}

	async function loadClasses() {
		classesLoading = true;
		try {
			const response = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes`);
			if (response.ok) {
				classes = await response.json();
			} else {
				console.error('Failed to load classes');
				classes = [];
			}
		} catch (error) {
			console.error('Error loading classes:', error);
			classes = [];
		} finally {
			classesLoading = false;
		}
	}

	async function handleLogout() {
		isLoading = true;
		try {
			// 백엔드에 로그아웃 요청
			await fetchWithAuth(API_ENDPOINTS.LOGOUT, {
				method: 'POST',
			});
		} catch (error) {
			console.error('Logout error:', error);
		} finally {
			// 프론트엔드에서도 쿠키 삭제
			logout();
			goto('/');
		}
	}

	// 학생 생성 폼 열기
	function openCreateForm() {
		showCreateForm = true;
		createFormData = {
			id: '',
			password: '',
			phone: '',
			class_id: null,
			role: 'student',
		};
	}

	// 학생 생성 폼 닫기
	function closeCreateForm() {
		showCreateForm = false;
		createFormData = {
			id: '',
			password: '',
			phone: '',
			class_id: null,
			role: 'student',
		};
	}

	// 학생 생성 요청
	async function createStudent() {
		if (!createFormData.id || !createFormData.password) {
			alert('이름과 비밀번호를 입력해주세요.');
			return;
		}

		isCreating = true;
		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/users`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(createFormData),
			});

			if (res.ok) {
				alert('학생이 성공적으로 생성되었습니다.');
				closeCreateForm();
				await loadStudents();
			} else {
				const errorData = await res.json();
				alert(`학생 생성 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('학생 생성 오류:', e);
			alert('학생 생성 중 오류가 발생했습니다.');
		} finally {
			isCreating = false;
		}
	}

	// 학생 수정 폼 열기
	function openEditForm(student) {
		selectedStudent = student;
		showEditForm = true;
		editFormData = {
			id: student.id,
			password: '',
			phone: student.phone || '',
			class_id: student.class_id,
			role: student.role || 'student',
		};
	}

	// 학생 수정 폼 닫기
	function closeEditForm() {
		showEditForm = false;
		selectedStudent = null;
		editFormData = {
			id: '',
			password: '',
			phone: '',
			class_id: null,
			role: 'student',
		};
	}

	// 학생 수정 요청
	async function updateStudent() {
		if (!editFormData.id || !editFormData.password) {
			alert('학번과 비밀번호를 입력해주세요.');
			return;
		}

		isEditing = true;
		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/users/${selectedStudent.id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(editFormData),
			});

			if (res.ok) {
				alert('학생 정보가 성공적으로 수정되었습니다.');
				closeEditForm();
				await loadStudents();
			} else {
				const errorData = await res.json();
				alert(`학생 수정 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('학생 수정 오류:', e);
			alert('학생 수정 중 오류가 발생했습니다.');
		} finally {
			isEditing = false;
		}
	}

	// 학생 삭제 요청
	async function deleteStudent(student) {
		if (!confirm(`정말로 ${student.id} 학생을 삭제하시겠습니까?`)) {
			return;
		}

		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/users/${student.id}`, {
				method: 'DELETE',
			});

			if (res.ok) {
				alert('학생이 성공적으로 삭제되었습니다.');
				await loadStudents();
			} else {
				const errorData = await res.json();
				alert(`학생 삭제 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('학생 삭제 오류:', e);
			alert('학생 삭제 중 오류가 발생했습니다.');
		}
	}

	// 반 생성 폼 열기
	function openCreateClassForm() {
		showCreateClassForm = true;
		createClassFormData = {
			class_name: '',
		};
	}

	// 반 생성 폼 닫기
	function closeCreateClassForm() {
		showCreateClassForm = false;
		createClassFormData = {
			class_name: '',
		};
	}

	// 반 생성 요청
	async function createClass() {
		if (!createClassFormData.class_name) {
			alert('반 이름을 입력해주세요.');
			return;
		}

		isCreatingClass = true;
		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(createClassFormData),
			});

			if (res.ok) {
				alert('반이 성공적으로 생성되었습니다.');
				closeCreateClassForm();
				await loadClasses();
			} else {
				const errorData = await res.json();
				alert(`반 생성 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('반 생성 오류:', e);
			alert('반 생성 중 오류가 발생했습니다.');
		} finally {
			isCreatingClass = false;
		}
	}

	// 반 수정 폼 열기
	function openEditClassForm(classItem) {
		selectedClass = classItem;
		showEditClassForm = true;
		editClassFormData = {
			class_id: classItem.class_id,
			class_name: classItem.class_name,
		};
	}

	// 반 수정 폼 닫기
	function closeEditClassForm() {
		showEditClassForm = false;
		selectedClass = null;
		editClassFormData = {
			class_id: null,
			class_name: '',
		};
	}

	// 반 수정 요청
	async function updateClass() {
		if (!editClassFormData.class_name) {
			alert('반 이름을 입력해주세요.');
			return;
		}

		isEditingClass = true;
		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes/${selectedClass.class_id}`, {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(editClassFormData),
			});

			if (res.ok) {
				alert('반 정보가 성공적으로 수정되었습니다.');
				closeEditClassForm();
				await loadClasses();
			} else {
				const errorData = await res.json();
				alert(`반 수정 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('반 수정 오류:', e);
			alert('반 수정 중 오류가 발생했습니다.');
		} finally {
			isEditingClass = false;
		}
	}

	// 반 삭제 요청
	async function deleteClass(classItem) {
		if (!confirm(`정말로 "${classItem.class_name}" 반을 삭제하시겠습니까?`)) {
			return;
		}

		try {
			const res = await fetchWithAuth(`${API_ENDPOINTS.API_BASE_URL}/classes/${classItem.class_id}`, {
				method: 'DELETE',
			});

			if (res.ok) {
				alert('반이 성공적으로 삭제되었습니다.');
				await loadClasses();
			} else {
				const errorData = await res.json();
				alert(`반 삭제 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('반 삭제 오류:', e);
			alert('반 삭제 중 오류가 발생했습니다.');
		}
	}

	// 반 이름으로 반 ID 찾기
	function getClassNameById(classId) {
		if (!classId) return '-';
		const classItem = classes.find(c => c.class_id === classId);
		return classItem ? classItem.class_name : classId;
	}

	// 반별로 학생 그룹화
	function getStudentsByClass() {
		const grouped = {};
		
		// 반이 없는 학생들을 위한 그룹
		grouped['no-class'] = {
			class_id: null,
			class_name: '반 미지정',
			students: []
		};
		
		// 각 반별로 학생 그룹화
		classes.forEach(classItem => {
			grouped[classItem.class_id] = {
				class_id: classItem.class_id,
				class_name: classItem.class_name,
				students: []
			};
		});
		
		// 학생들을 해당 반에 배치
		students.forEach(student => {
			if (student.class_id && grouped[student.class_id]) {
				grouped[student.class_id].students.push(student);
			} else {
				grouped['no-class'].students.push(student);
			}
		});
		
		return Object.values(grouped).filter(group => group.students.length > 0);
	}

	// 반 토글 함수
	function toggleClass(classId) {
		if (expandedClasses.has(classId)) {
			expandedClasses.delete(classId);
		} else {
			expandedClasses.add(classId);
		}
		expandedClasses = expandedClasses; // Svelte 반응성 트리거
	}

	// 시험 기간 관련 함수들
	async function loadExamPeriods(classId) {
		examPeriodsLoading = true;
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIODS_BY_CLASS(classId));
			if (response.ok) {
				examPeriods = await response.json();
			} else {
				console.error('Failed to load exam periods');
				examPeriods = [];
			}
		} catch (error) {
			console.error('Error loading exam periods:', error);
			examPeriods = [];
		} finally {
			examPeriodsLoading = false;
		}
	}

	// 전체 시험 기간 로드 (모든 반의 시험 기간)
	async function loadAllExamPeriods() {
		try {
			const response = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIODS);
			if (response.ok) {
				const periods = await response.json();
				// 각 시험 기간에 반 정보 추가
				periods.forEach(period => {
					const classItem = classes.find(c => c.class_id === period.class_id);
					period.class_name = classItem ? classItem.class_name : '알 수 없는 반';
				});
				allExamPeriods = periods;
				calculateExamPeriodsSummary();
			} else {
				console.error('Failed to load all exam periods');
				allExamPeriods = [];
			}
		} catch (error) {
			console.error('Error loading all exam periods:', error);
			allExamPeriods = [];
		}
	}

	// 시험 기간 요약 정보 계산
	function calculateExamPeriodsSummary() {
		const now = new Date();
		const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
		
		examPeriodsSummary = {
			total: allExamPeriods.length,
			upcoming: 0,
			ongoing: 0,
			completed: 0
		};

		allExamPeriods.forEach(period => {
			const startDate = new Date(period.start_date);
			const endDate = new Date(period.end_date);
			
			if (today < startDate) {
				examPeriodsSummary.upcoming++;
			} else if (today >= startDate && today <= endDate) {
				examPeriodsSummary.ongoing++;
			} else if (today > endDate) {
				examPeriodsSummary.completed++;
			}
		});
	}

	// 반 클릭 시 시험 기간 표시
	async function showExamPeriodsForClass(classItem) {
		selectedClassForExam = classItem;
		showExamPeriods = true;
		await loadExamPeriods(classItem.class_id);
	}

	// 전체 시험 기간 관리 화면 표시
	async function showAllExamPeriods() {
		selectedClassForExam = null;
		showExamPeriods = true;
		await loadAllExamPeriods();
	}

	// 시험 기간 생성 폼 열기
	function openCreateExamPeriodForm(classItem) {
		showCreateExamPeriodForm = true;
		createExamPeriodFormData = {
			class_id: classItem.class_id,
			name: '',
			description: '',
			start_date: '',
			end_date: '',
		};
	}

	// 시험 기간 생성 폼 닫기
	function closeCreateExamPeriodForm() {
		showCreateExamPeriodForm = false;
		createExamPeriodFormData = {
			class_id: null,
			name: '',
			description: '',
			start_date: '',
			end_date: '',
			english_date: '',
		};
	}

	// 시험 기간 생성 요청
	async function createExamPeriod() {
		if (!createExamPeriodFormData.name || !createExamPeriodFormData.start_date || !createExamPeriodFormData.end_date) {
			alert('시험 기간명, 시작일, 종료일을 모두 입력해주세요.');
			return;
		}

		isCreatingExamPeriod = true;
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIODS, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(createExamPeriodFormData),
			});

			if (res.ok) {
				alert('시험 기간이 성공적으로 생성되었습니다.');
				closeCreateExamPeriodForm();
				await loadExamPeriods(createExamPeriodFormData.class_id);
				await loadAllExamPeriods(); // 전체 시험 기간 요약 업데이트
			} else {
				const errorData = await res.json();
				alert(`시험 기간 생성 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('시험 기간 생성 오류:', e);
			alert('시험 기간 생성 중 오류가 발생했습니다.');
		} finally {
			isCreatingExamPeriod = false;
		}
	}

	// 시험 기간 수정 폼 열기
	function openEditExamPeriodForm(examPeriod) {
		selectedExamPeriod = examPeriod;
		showEditExamPeriodForm = true;
		editExamPeriodFormData = {
			exam_period_id: examPeriod.exam_period_id,
			class_id: examPeriod.class_id,
			name: examPeriod.name,
			description: examPeriod.description || '',
			start_date: examPeriod.start_date,
			end_date: examPeriod.end_date,
			english_date: examPeriod.english_date || '',
		};
	}

	// 시험 기간 수정 폼 닫기
	function closeEditExamPeriodForm() {
		showEditExamPeriodForm = false;
		selectedExamPeriod = null;
		editExamPeriodFormData = {
			exam_period_id: null,
			class_id: null,
			name: '',
			description: '',
			start_date: '',
			end_date: '',
			english_date: '',
		};
	}

	// 시험 기간 수정 요청
	async function updateExamPeriod() {
		if (!editExamPeriodFormData.name || !editExamPeriodFormData.start_date || !editExamPeriodFormData.end_date) {
			alert('시험 기간명, 시작일, 종료일을 모두 입력해주세요.');
			return;
		}

		isEditingExamPeriod = true;
		try {
			const res = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIOD_BY_ID(selectedExamPeriod.exam_period_id), {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(editExamPeriodFormData),
			});

			if (res.ok) {
				alert('시험 기간이 성공적으로 수정되었습니다.');
				closeEditExamPeriodForm();
				await loadExamPeriods(editExamPeriodFormData.class_id);
				await loadAllExamPeriods(); // 전체 시험 기간 요약 업데이트
			} else {
				const errorData = await res.json();
				alert(`시험 기간 수정 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('시험 기간 수정 오류:', e);
			alert('시험 기간 수정 중 오류가 발생했습니다.');
		} finally {
			isEditingExamPeriod = false;
		}
	}

	// 시험 기간 삭제 요청
	async function deleteExamPeriod(examPeriod) {
		if (!confirm(`정말로 "${examPeriod.name}" 시험 기간을 삭제하시겠습니까?`)) {
			return;
		}

		try {
			const res = await fetchWithAuth(API_ENDPOINTS.EXAM_PERIOD_BY_ID(examPeriod.exam_period_id), {
				method: 'DELETE',
			});

			if (res.ok) {
				alert('시험 기간이 성공적으로 삭제되었습니다.');
				await loadExamPeriods(examPeriod.class_id);
				await loadAllExamPeriods(); // 전체 시험 기간 요약 업데이트
			} else {
				const errorData = await res.json();
				alert(`시험 기간 삭제 실패: ${errorData.error || '알 수 없는 오류가 발생했습니다.'}`);
			}
		} catch (e) {
			console.error('시험 기간 삭제 오류:', e);
			alert('시험 기간 삭제 중 오류가 발생했습니다.');
		}
	}
</script>

<svelte:head>
	<title>Julia - 관리자 대시보드</title>
	<meta name="description" content="Julia 시스템 관리자 대시보드" />
</svelte:head>

<div class="dashboard-container">
	<header class="dashboard-header">
		<div class="header-left">
			<h1>관리자 대시보드</h1>
			<div class="user-info">
				<span>안녕하세요, {userLoading ? '로딩 중...' : username}님!</span>
				<span class="user-role admin">관리자</span>
			</div>
		</div>
		<div class="header-right">
			<button class="btn ghost" on:click={() => goto('/calendar')}>
				<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
					<path
						d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z"
						stroke="currentColor"
						fill="none"
						stroke-width="2"
					/>
				</svg>
				캘린더
			</button>
			{#if userRole === 'admin'}
				<button class="btn ghost" on:click={() => goto('/upload')}>
					<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
						<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" stroke="currentColor" fill="none" stroke-width="2"/>
						<polyline points="14,2 14,8 20,8" stroke="currentColor" fill="none" stroke-width="2"/>
						<line x1="16" y1="13" x2="8" y2="13" stroke="currentColor" stroke-width="2"/>
						<line x1="16" y1="17" x2="8" y2="17" stroke="currentColor" stroke-width="2"/>
						<polyline points="10,9 9,9 8,9" stroke="currentColor" fill="none" stroke-width="2"/>
					</svg>
					파일 업로드
				</button>
				<button class="btn primary" on:click={showAllExamPeriods}>
					<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
						<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2" />
					</svg>
					시험 기간 관리
				</button>
			{/if}
			<button class="logout-button" on:click={handleLogout} disabled={isLoading}>
				{isLoading ? '로그아웃 중...' : '로그아웃'}
			</button>
		</div>
	</header>

	<main class="dashboard-content">
		{#if userRole !== 'admin'}
			<div class="access-denied">
				<h2>접근 권한이 없습니다</h2>
				<p>관리자만 이 페이지에 접근할 수 있습니다.</p>
				<button class="btn primary" on:click={() => goto('/calendar')}> 캘린더로 돌아가기 </button>
			</div>
		{:else}
			<!-- 시험 기간 요약 섹션 -->
			<div class="summary-section">
				<div class="summary-header">
					<h2>시험 기간 현황</h2>
					<button class="btn primary" on:click={showAllExamPeriods}>
						<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z" stroke="currentColor" fill="none" stroke-width="2" />
						</svg>
						전체 관리
					</button>
				</div>
				<div class="summary-cards">
					<div class="summary-card total">
						<div class="card-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M9 12l2 2 4-4M21 12c0 4.97-4.03 9-9 9s-9-4.03-9-9 4.03-9 9-9 9 4.03 9 9z" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
						</div>
						<div class="card-content">
							<div class="card-number">{examPeriodsSummary.total}</div>
							<div class="card-label">전체 시험 기간</div>
						</div>
					</div>
					<div class="summary-card upcoming">
						<div class="card-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
						</div>
						<div class="card-content">
							<div class="card-number">{examPeriodsSummary.upcoming}</div>
							<div class="card-label">예정된 시험</div>
						</div>
					</div>
					<div class="summary-card ongoing">
						<div class="card-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
						</div>
						<div class="card-content">
							<div class="card-number">{examPeriodsSummary.ongoing}</div>
							<div class="card-label">진행 중인 시험</div>
						</div>
					</div>
					<div class="summary-card completed">
						<div class="card-icon">
							<svg width="24" height="24" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
						</div>
						<div class="card-content">
							<div class="card-number">{examPeriodsSummary.completed}</div>
							<div class="card-label">완료된 시험</div>
						</div>
					</div>
				</div>
			</div>

			<div class="dashboard-grid">
				<!-- 반 관리 섹션 -->
				<div class="classes-section">
					<div class="section-header">
						<h2>반 관리</h2>
						<button class="btn primary" on:click={openCreateClassForm}>
							<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
							새 반 추가
						</button>
					</div>

					<div class="classes-list">
						{#if classesLoading}
							<div class="loading">
								<div class="skeleton head" />
								<div class="skeleton row" />
								<div class="skeleton row" />
							</div>
						{:else if classes.length === 0}
							<div class="empty-state">
								<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
									<path
										d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
										stroke="currentColor"
										fill="none"
										stroke-width="2"
									/>
								</svg>
								<p>등록된 반이 없습니다.</p>
								<button class="btn primary" on:click={openCreateClassForm}> 첫 번째 반 추가하기 </button>
							</div>
						{:else}
							<div class="classes-table">
								<table>
									<thead>
										<tr>
											<th>반 번호</th>
											<th>반 이름</th>
											<th>작업</th>
										</tr>
									</thead>
									<tbody>
										{#each classes as classItem}
											<tr class="clickable-row" on:click={() => showExamPeriodsForClass(classItem)}>
												<td>{classItem.class_id}</td>
												<td>{classItem.class_name}</td>
												<td>
													<div class="actions">
														<button
															class="btn-icon edit"
															on:click|stopPropagation={() => openEditClassForm(classItem)}
															title="수정"
														>
															<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
																<path
																	d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
																<path
																	d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
															</svg>
														</button>
														<button
															class="btn-icon delete"
															on:click|stopPropagation={() => deleteClass(classItem)}
															title="삭제"
														>
															<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
																<path
																	d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0v14M10 11v6M14 11v6"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
															</svg>
														</button>
													</div>
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{/if}
					</div>
				</div>

				<!-- 학생 관리 섹션 -->
				<div class="students-section">
					<div class="section-header">
						<h2>학생 관리</h2>
						<button class="btn primary" on:click={openCreateForm}>
							<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
								<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2" />
							</svg>
							새 학생 추가
						</button>
					</div>

					<!-- 반별 학생 목록 -->
					<div class="students-list">
						<h3>반별 학생 목록</h3>
						{#if studentsLoading || classesLoading}
							<div class="loading">
								<div class="skeleton head" />
								<div class="skeleton row" />
								<div class="skeleton row" />
								<div class="skeleton row" />
							</div>
						{:else if students.length === 0}
							<div class="empty-state">
								<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
									<path
										d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2M12 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8z"
										stroke="currentColor"
										fill="none"
										stroke-width="2"
									/>
								</svg>
								<p>등록된 학생이 없습니다.</p>
								<button class="btn primary" on:click={openCreateForm}> 첫 번째 학생 추가하기 </button>
							</div>
						{:else}
							<div class="class-groups">
								{#each getStudentsByClass() as classGroup}
									<div class="class-group">
										<div class="class-header" on:click={() => toggleClass(classGroup.class_id)}>
											<div class="class-info">
												<span class="class-name">{classGroup.class_name}</span>
												<span class="student-count">({classGroup.students.length}명)</span>
											</div>
											<button class="toggle-button" class:expanded={expandedClasses.has(classGroup.class_id)}>
												<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
													<path d="M6 9l6 6 6-6" stroke="currentColor" fill="none" stroke-width="2" />
												</svg>
											</button>
										</div>
										
										<div class="students-table" class:expanded={expandedClasses.has(classGroup.class_id)}>
											<table>
													<thead>
														<tr>
															<th>이름</th>
															<th>전화번호</th>
															<th>역할</th>
															<th>가입일</th>
															<th>작업</th>
														</tr>
													</thead>
													<tbody>
														{#each classGroup.students as student}
															<tr>
																<td>{student.id}</td>
																<td>{student.phone || '-'}</td>
																<td>
																	<span class="role-badge {student.role}">
																		{student.role === 'admin' ? '관리자' : '학생'}
																	</span>
																</td>
																<td>{new Date(student.created_at).toLocaleDateString()}</td>
																<td>
																	<div class="actions">
																		<button
																			class="btn-icon edit"
																			on:click={() => openEditForm(student)}
																			title="수정"
																		>
																			<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
																<path
																	d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
																<path
																	d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
															</svg>
														</button>
														<button
															class="btn-icon delete"
															on:click={() => deleteStudent(student)}
															title="삭제"
														>
															<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
																<path
																	d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0v14M10 11v6M14 11v6"
																	stroke="currentColor"
																	fill="none"
																	stroke-width="2"
																/>
															</svg>
														</button>
													</div>
												</td>
											</tr>
										{/each}
																						</tbody>
												</table>
											</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			</div>
		{/if}
	</main>
</div>

<!-- 학생 생성 모달 -->
{#if showCreateForm}
	<div
		class="modal-overlay"
		on:click={closeCreateForm}
		on:keydown={(e) => e.key === 'Escape' && closeCreateForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>새 학생 추가</h3>
			<button class="btn-close" on:click={closeCreateForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={createStudent}>
			<div class="form-group">
				<label for="create-id">이름 *</label>
				<input
					id="create-id"
					type="text"
					bind:value={createFormData.id}
					placeholder="이름을 입력하세요"
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-password">비밀번호 *</label>
				<input
					id="create-password"
					type="password"
					bind:value={createFormData.password}
					placeholder="비밀번호를 입력하세요"
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-phone">전화번호</label>
				<input
					id="create-phone"
					type="tel"
					bind:value={createFormData.phone}
					placeholder="전화번호를 입력하세요"
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-class">반</label>
				<select id="create-class" bind:value={createFormData.class_id} class="form-input">
					<option value={null}>반을 선택하세요</option>
					{#each classes as classItem}
						<option value={classItem.class_id}>{classItem.class_name}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="create-role">역할</label>
				<select id="create-role" bind:value={createFormData.role} class="form-input">
					<option value="student">학생</option>
					<option value="admin">관리자</option>
				</select>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeCreateForm} disabled={isCreating}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isCreating}>
					{isCreating ? '생성 중...' : '학생 생성'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 학생 수정 모달 -->
{#if showEditForm}
	<div
		class="modal-overlay"
		on:click={closeEditForm}
		on:keydown={(e) => e.key === 'Escape' && closeEditForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>학생 정보 수정</h3>
			<button class="btn-close" on:click={closeEditForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={updateStudent}>
			<div class="form-group">
				<label for="edit-id">이름 *</label>
				<input
					id="edit-id"
					type="text"
					bind:value={editFormData.id}
					placeholder="이름을 입력하세요"
					required
					class="form-input"
					readonly
				/>
			</div>

			<div class="form-group">
				<label for="edit-password">새 비밀번호 *</label>
				<input
					id="edit-password"
					type="password"
					bind:value={editFormData.password}
					placeholder="새 비밀번호를 입력하세요"
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="edit-phone">전화번호</label>
				<input
					id="edit-phone"
					type="tel"
					bind:value={editFormData.phone}
					placeholder="전화번호를 입력하세요"
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="edit-class">반</label>
				<select id="edit-class" bind:value={editFormData.class_id} class="form-input">
					<option value={null}>반을 선택하세요</option>
					{#each classes as classItem}
						<option value={classItem.class_id}>{classItem.class_name}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="edit-role">역할</label>
				<select id="edit-role" bind:value={editFormData.role} class="form-input">
					<option value="student">학생</option>
					<option value="admin">관리자</option>
				</select>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeEditForm} disabled={isEditing}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isEditing}>
					{isEditing ? '수정 중...' : '학생 수정'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 반 생성 모달 -->
{#if showCreateClassForm}
	<div
		class="modal-overlay"
		on:click={closeCreateClassForm}
		on:keydown={(e) => e.key === 'Escape' && closeCreateClassForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>새 반 추가</h3>
			<button class="btn-close" on:click={closeCreateClassForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={createClass}>
			<div class="form-group">
				<label for="create-class-name">반 이름 *</label>
				<input
					id="create-class-name"
					type="text"
					bind:value={createClassFormData.class_name}
					placeholder="반 이름을 입력하세요 (예: 1학년 1반)"
					required
					class="form-input"
				/>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeCreateClassForm} disabled={isCreatingClass}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isCreatingClass}>
					{isCreatingClass ? '생성 중...' : '반 생성'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 반 수정 모달 -->
{#if showEditClassForm}
	<div
		class="modal-overlay"
		on:click={closeEditClassForm}
		on:keydown={(e) => e.key === 'Escape' && closeEditClassForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>반 정보 수정</h3>
			<button class="btn-close" on:click={closeEditClassForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={updateClass}>
			<div class="form-group">
				<label for="edit-class-name">반 이름 *</label>
				<input
					id="edit-class-name"
					type="text"
					bind:value={editClassFormData.class_name}
					placeholder="반 이름을 입력하세요 (예: 1학년 1반)"
					required
					class="form-input"
				/>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeEditClassForm} disabled={isEditingClass}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isEditingClass}>
					{isEditingClass ? '수정 중...' : '반 수정'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 시험 기간 목록 모달 -->
{#if showExamPeriods}
	<div
		class="modal-overlay"
		on:click={() => showExamPeriods = false}
		on:keydown={(e) => e.key === 'Escape' && (showExamPeriods = false)}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content exam-periods-modal" on:click|stopPropagation>
		<div class="modal-header">
			<h3>
				{selectedClassForExam ? `${selectedClassForExam.class_name} 시험 기간` : '전체 시험 기간 관리'}
			</h3>
			<div class="modal-header-actions">
				{#if selectedClassForExam}
					<button class="btn primary" on:click={() => openCreateExamPeriodForm(selectedClassForExam)}>
						<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M12 5v14M5 12h14" stroke="currentColor" fill="none" stroke-width="2" />
						</svg>
						시험 기간 추가
					</button>
				{:else}
					<button class="btn ghost" on:click={() => showExamPeriods = false}>
						<svg width="16" height="16" viewBox="0 0 24 24" aria-hidden="true">
							<path d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" stroke="currentColor" fill="none" stroke-width="2" />
						</svg>
						반별 관리
					</button>
				{/if}
				<button class="btn-close" on:click={() => showExamPeriods = false} aria-label="닫기">
					<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
						<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
					</svg>
				</button>
			</div>
		</div>

		<div class="exam-periods-content">
			{#if !selectedClassForExam}
				<div class="all-exam-periods-list">
					{#if examPeriodsLoading}
						<div class="loading">
							<div class="skeleton head" />
							<div class="skeleton row" />
							<div class="skeleton row" />
						</div>
					{:else if allExamPeriods.length === 0}
						<div class="empty-state">
							<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
								<path
									d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z"
									stroke="currentColor"
									fill="none"
									stroke-width="2"
								/>
							</svg>
							<p>등록된 시험 기간이 없습니다.</p>
						</div>
					{:else}
						<div class="all-exam-periods-table">
							<table>
								<thead>
									<tr>
										<th>반</th>
										<th>시험 기간명</th>
										<th>설명</th>
										<th>시작일</th>
										<th>종료일</th>
										<th>영어 시험</th>
										<th>작업</th>
									</tr>
								</thead>
								<tbody>
									{#each allExamPeriods as examPeriod}
										<tr>
											<td>
												<button class="class-link" on:click={() => showExamPeriodsForClass({class_id: examPeriod.class_id, class_name: examPeriod.class_name})}>
													{examPeriod.class_name}
												</button>
											</td>
											<td>{examPeriod.name}</td>
											<td>{examPeriod.description || '-'}</td>
											<td>{new Date(examPeriod.start_date).toLocaleDateString()}</td>
											<td>{new Date(examPeriod.end_date).toLocaleDateString()}</td>
											<td>
												{#if examPeriod.english_date}
													<span class="english-date-badge">
														📚 {new Date(examPeriod.english_date).toLocaleDateString()}
													</span>
												{:else}
													<span class="no-english-date">-</span>
												{/if}
											</td>
											<td>
												<div class="actions">
													<button
														class="btn-icon edit"
														on:click={() => openEditExamPeriodForm(examPeriod)}
														title="수정"
													>
														<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
															<path
																d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
															<path
																d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
														</svg>
													</button>
													<button
														class="btn-icon delete"
														on:click={() => deleteExamPeriod(examPeriod)}
														title="삭제"
													>
														<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
															<path
																d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0v14M10 11v6M14 11v6"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
														</svg>
													</button>
												</div>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			{:else}
				<div class="exam-periods-list">
					{#if examPeriodsLoading}
						<div class="loading">
							<div class="skeleton head" />
							<div class="skeleton row" />
							<div class="skeleton row" />
						</div>
					{:else if examPeriods.length === 0}
						<div class="empty-state">
							<svg width="48" height="48" viewBox="0 0 24 24" aria-hidden="true">
								<path
									d="M8 2v4M16 2v4M3 10h18M5 4h14a2 2 0 012 2v14a2 2 0 01-2 2H5a2 2 0 01-2-2V6a2 2 0 012-2z"
									stroke="currentColor"
									fill="none"
									stroke-width="2"
								/>
							</svg>
							<p>등록된 시험 기간이 없습니다.</p>
							<button class="btn primary" on:click={() => openCreateExamPeriodForm(selectedClassForExam)}>
								첫 번째 시험 기간 추가하기
							</button>
						</div>
					{:else}
						<div class="exam-periods-table">
							<table>
								<thead>
									<tr>
										<th>시험 기간명</th>
										<th>설명</th>
										<th>시작일</th>
										<th>종료일</th>
										<th>영어 시험</th>
										<th>작업</th>
									</tr>
								</thead>
								<tbody>
									{#each examPeriods as examPeriod}
										<tr>
											<td>{examPeriod.name}</td>
											<td>{examPeriod.description || '-'}</td>
											<td>{new Date(examPeriod.start_date).toLocaleDateString()}</td>
											<td>{new Date(examPeriod.end_date).toLocaleDateString()}</td>
											<td>
												{#if examPeriod.english_date}
													<span class="english-date-badge">
														📚 {new Date(examPeriod.english_date).toLocaleDateString()}
													</span>
												{:else}
													<span class="no-english-date">-</span>
												{/if}
											</td>
											<td>
												<div class="actions">
													<button
														class="btn-icon edit"
														on:click={() => openEditExamPeriodForm(examPeriod)}
														title="수정"
													>
														<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
															<path
																d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
															<path
																d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
														</svg>
													</button>
													<button
														class="btn-icon delete"
														on:click={() => deleteExamPeriod(examPeriod)}
														title="삭제"
													>
														<svg width="14" height="14" viewBox="0 0 24 24" aria-hidden="true">
															<path
																d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0v14M10 11v6M14 11v6"
																stroke="currentColor"
																fill="none"
																stroke-width="2"
															/>
														</svg>
													</button>
												</div>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
{/if}

<!-- 시험 기간 생성 모달 -->
{#if showCreateExamPeriodForm}
	<div
		class="modal-overlay"
		on:click={closeCreateExamPeriodForm}
		on:keydown={(e) => e.key === 'Escape' && closeCreateExamPeriodForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>새 시험 기간 추가</h3>
			<button class="btn-close" on:click={closeCreateExamPeriodForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={createExamPeriod}>
			<div class="form-group">
				<label for="create-exam-name">시험 기간명 *</label>
				<input
					id="create-exam-name"
					type="text"
					bind:value={createExamPeriodFormData.name}
					placeholder="예: 1학기 중간고사"
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-exam-description">설명</label>
				<textarea
					id="create-exam-description"
					bind:value={createExamPeriodFormData.description}
					placeholder="시험 기간에 대한 설명을 입력하세요"
					class="form-input"
					rows="3"
				></textarea>
			</div>

			<div class="form-group">
				<label for="create-start-date">시작일 *</label>
				<input
					id="create-start-date"
					type="date"
					bind:value={createExamPeriodFormData.start_date}
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-end-date">종료일 *</label>
				<input
					id="create-end-date"
					type="date"
					bind:value={createExamPeriodFormData.end_date}
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="create-english-date">영어 시험 날짜 (선택사항)</label>
				<input
					id="create-english-date"
					type="date"
					bind:value={createExamPeriodFormData.english_date}
					class="form-input"
				/>
				<small class="form-help">영어 시험이 별도 날짜에 있는 경우 입력하세요</small>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeCreateExamPeriodForm} disabled={isCreatingExamPeriod}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isCreatingExamPeriod}>
					{isCreatingExamPeriod ? '생성 중...' : '시험 기간 생성'}
				</button>
			</div>
		</form>
	</div>
{/if}

<!-- 시험 기간 수정 모달 -->
{#if showEditExamPeriodForm}
	<div
		class="modal-overlay"
		on:click={closeEditExamPeriodForm}
		on:keydown={(e) => e.key === 'Escape' && closeEditExamPeriodForm()}
		role="button"
		tabindex="0"
		aria-label="모달 닫기"
	/>
	<div class="modal-content" on:click|stopPropagation>
		<div class="modal-header">
			<h3>시험 기간 수정</h3>
			<button class="btn-close" on:click={closeEditExamPeriodForm} aria-label="닫기">
				<svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">
					<path d="M18 6L6 18M6 6l12 12" stroke="currentColor" fill="none" stroke-width="2" />
				</svg>
			</button>
		</div>

		<form on:submit|preventDefault={updateExamPeriod}>
			<div class="form-group">
				<label for="edit-exam-name">시험 기간명 *</label>
				<input
					id="edit-exam-name"
					type="text"
					bind:value={editExamPeriodFormData.name}
					placeholder="예: 1학기 중간고사"
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="edit-exam-description">설명</label>
				<textarea
					id="edit-exam-description"
					bind:value={editExamPeriodFormData.description}
					placeholder="시험 기간에 대한 설명을 입력하세요"
					class="form-input"
					rows="3"
				></textarea>
			</div>

			<div class="form-group">
				<label for="edit-start-date">시작일 *</label>
				<input
					id="edit-start-date"
					type="date"
					bind:value={editExamPeriodFormData.start_date}
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="edit-end-date">종료일 *</label>
				<input
					id="edit-end-date"
					type="date"
					bind:value={editExamPeriodFormData.end_date}
					required
					class="form-input"
				/>
			</div>

			<div class="form-group">
				<label for="edit-english-date">영어 시험 날짜 (선택사항)</label>
				<input
					id="edit-english-date"
					type="date"
					bind:value={editExamPeriodFormData.english_date}
					class="form-input"
				/>
				<small class="form-help">영어 시험이 별도 날짜에 있는 경우 입력하세요</small>
			</div>

			<div class="form-actions">
				<button type="button" class="btn ghost" on:click={closeEditExamPeriodForm} disabled={isEditingExamPeriod}>
					취소
				</button>
				<button type="submit" class="btn primary" disabled={isEditingExamPeriod}>
					{isEditingExamPeriod ? '수정 중...' : '시험 기간 수정'}
				</button>
			</div>
		</form>
	</div>
{/if}

<style>
	:root {
		--bg: #f6f7fb;
		--card: #ffffff;
		--text: #1f2937;
		--muted: #6b7280;
		--line: #e5e7eb;
		--brand: #3b82f6;
		--brand-600: #2563eb;
		--success: #10b981;
		--warning: #f59e0b;
		--danger: #ef4444;
	}



	* {
		box-sizing: border-box;
	}

	.dashboard-container {
		min-height: 100vh;
		background: var(--bg);
		color: var(--text);
	}

	.dashboard-header {
		background: var(--card);
		padding: 1rem 2rem;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		display: flex;
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid var(--line);
	}

	.header-left h1 {
		color: var(--text);
		font-size: 1.5rem;
		font-weight: 600;
		margin: 0 0 0.5rem 0;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.user-info span {
		color: var(--muted);
		font-size: 0.875rem;
	}

	.user-role {
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 600;
	}

	.user-role.admin {
		background: var(--brand);
		color: white;
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.btn {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		border-radius: 0.5rem;
		padding: 0.5rem 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
	}

	.btn:hover {
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	}

	.btn.ghost {
		background: transparent;
		border-color: var(--line);
	}

	.btn.primary {
		background: var(--brand);
		border-color: var(--brand);
		color: white;
	}

	.btn svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.logout-button {
		padding: 0.5rem 1rem;
		background: var(--danger);
		color: white;
		border: none;
		border-radius: 0.375rem;
		cursor: pointer;
		font-size: 0.875rem;
		transition: background-color 0.2s ease;
	}

	.logout-button:hover {
		background: #dc2626;
	}

	.dashboard-content {
		padding: 2rem;
		max-width: 1200px;
		margin: 0 auto;
		/* app.css의 main 스타일 오버라이드 */
		display: block !important;
		align-items: unset !important;
		justify-content: unset !important;
		min-height: unset !important;
		padding-top: 2rem !important;
		padding-right: 2rem !important;
	}

	/* 요약 섹션 스타일 */
	.summary-section {
		background: var(--card);
		border-radius: 0.75rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		overflow: hidden;
		margin-bottom: 2rem;
	}

	.summary-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem 2rem;
		border-bottom: 1px solid var(--line);
	}

	.summary-header h2 {
		color: var(--text);
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.summary-cards {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 1rem;
		padding: 1.5rem 2rem;
	}

	.summary-card {
		background: var(--bg);
		border-radius: 0.5rem;
		padding: 1.5rem;
		display: flex;
		align-items: center;
		gap: 1rem;
		transition: transform 0.2s ease, box-shadow 0.2s ease;
	}

	.summary-card:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
	}

	.card-icon {
		width: 48px;
		height: 48px;
		border-radius: 0.5rem;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.card-content {
		flex: 1;
	}

	.card-number {
		font-size: 2rem;
		font-weight: 700;
		line-height: 1;
		margin-bottom: 0.25rem;
	}

	.card-label {
		font-size: 0.875rem;
		color: var(--muted);
		font-weight: 500;
	}

	/* 카드 색상 테마 */
	.summary-card.total .card-icon {
		background: rgba(59, 130, 246, 0.1);
		color: var(--brand);
	}

	.summary-card.total .card-number {
		color: var(--brand);
	}

	.summary-card.upcoming .card-icon {
		background: rgba(245, 158, 11, 0.1);
		color: var(--warning);
	}

	.summary-card.upcoming .card-number {
		color: var(--warning);
	}

	.summary-card.ongoing .card-icon {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success);
	}

	.summary-card.ongoing .card-number {
		color: var(--success);
	}

	.summary-card.completed .card-icon {
		background: rgba(107, 114, 128, 0.1);
		color: var(--muted);
	}

	.summary-card.completed .card-number {
		color: var(--muted);
	}

	.dashboard-grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 2rem;
		align-items: start; /* 각 섹션이 독립적인 높이를 가지도록 */
	}

	.classes-section {
		background: var(--card);
		border-radius: 0.75rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		overflow: hidden;
		height: fit-content; /* 내용에 맞는 높이로 설정 */
	}

	.classes-list {
		padding: 1.5rem 2rem;
	}

	.classes-list h3 {
		color: var(--text);
		font-size: 1rem;
		font-weight: 600;
		margin: 0 0 1rem 0;
	}

	.classes-table {
		overflow-x: auto;
	}

	.class-groups {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.class-group {
		border: 1px solid var(--line);
		border-radius: 0.5rem;
		overflow: hidden;
	}

	.class-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 1.5rem;
		background: var(--bg);
		cursor: pointer;
		transition: background-color 0.2s ease;
		border-bottom: 1px solid var(--line);
	}

	.class-header:hover {
		background: #f0f0f0;
	}

	.class-info {
		display: flex;
		align-items: center;
		gap: 0.5rem;
	}

	.class-name {
		font-weight: 600;
		color: var(--text);
		font-size: 1rem;
	}

	.student-count {
		color: var(--muted);
		font-size: 0.875rem;
	}

	.toggle-button {
		background: none;
		border: none;
		padding: 0.5rem;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: transform 0.1s ease;
		color: var(--muted);
	}

	.toggle-button:hover {
		background: rgba(0, 0, 0, 0.05);
	}

	.toggle-button.expanded {
		transform: rotate(180deg);
	}

	.toggle-button svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	.students-table {
		max-height: 0;
		overflow: hidden;
		transition: max-height 0.3s ease-in-out;
	}

	.students-table.expanded {
		max-height: 400px; /* 적절한 높이로 제한 */
		overflow-y: auto; /* 세로 스크롤 활성화 */
		transition: max-height 0.3s ease-in-out;
	}

	/* 테이블이 펼쳐질 때의 스타일 */
	.students-table.expanded table {
		opacity: 1;
		transform: translateY(0);
		transition: opacity 0.3s ease-in-out, transform 0.3s ease-in-out;
	}

	/* 테이블이 접혀있을 때의 스타일 */
	.students-table table {
		opacity: 0;
		transform: translateY(-10px);
		transition: opacity 0.3s ease-in-out, transform 0.3s ease-in-out;
	}

	.students-table.expanded table {
		opacity: 1;
		transform: translateY(0);
	}

	/* 스크롤바 스타일링 */
	.students-table.expanded::-webkit-scrollbar {
		width: 8px;
	}

	.students-table.expanded::-webkit-scrollbar-track {
		background: var(--bg);
		border-radius: 4px;
	}

	.students-table.expanded::-webkit-scrollbar-thumb {
		background: var(--line);
		border-radius: 4px;
	}

	.students-table.expanded::-webkit-scrollbar-thumb:hover {
		background: var(--muted);
	}

	.access-denied {
		text-align: center;
		padding: 4rem 2rem;
		background: var(--card);
		border-radius: 0.75rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	}

	.access-denied h2 {
		color: var(--text);
		margin-bottom: 1rem;
		font-size: 1.5rem;
	}

	.access-denied p {
		color: var(--muted);
		margin-bottom: 2rem;
	}

	.students-section {
		background: var(--card);
		border-radius: 0.75rem;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		overflow: hidden;
		height: fit-content; /* 내용에 맞는 높이로 설정 */
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1.5rem 2rem;
		border-bottom: 1px solid var(--line);
	}

	.section-header h2 {
		color: var(--text);
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.students-list {
		padding: 1.5rem 2rem;
	}

	.students-list h3 {
		color: var(--text);
		font-size: 1rem;
		font-weight: 600;
		margin: 0 0 1rem 0;
	}

	.loading {
		padding: 2rem;
	}

	.skeleton {
		border-radius: 0.5rem;
		background: linear-gradient(
			90deg,
			rgba(0, 0, 0, 0.06),
			rgba(0, 0, 0, 0.12),
			rgba(0, 0, 0, 0.06)
		);
		background-size: 200% 100%;
		animation: shine 1.2s linear infinite;
	}

	@keyframes shine {
		to {
			background-position: -200% 0;
		}
	}

	.skeleton.head {
		height: 1.5rem;
		margin-bottom: 1rem;
	}

	.skeleton.row {
		height: 3rem;
		margin: 0.5rem 0;
	}

	.empty-state {
		text-align: center;
		padding: 4rem 2rem;
		color: var(--muted);
	}

	.empty-state svg {
		opacity: 0.5;
		margin-bottom: 1rem;
	}

	.empty-state p {
		margin-bottom: 1.5rem;
		font-size: 1.125rem;
	}

	.students-table {
		overflow-x: auto;
	}

	table {
		width: 100%;
		border-collapse: collapse;
	}

	th,
	td {
		padding: 1rem;
		text-align: left;
		border-bottom: 1px solid var(--line);
		white-space: nowrap;
	}

	th {
		background: var(--bg);
		font-weight: 600;
		color: var(--text);
		font-size: 0.875rem;
	}

	td {
		color: var(--text);
		font-size: 0.875rem;
	}

	.role-badge {
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 600;
		white-space: nowrap;
	}

	.role-badge.student {
		background: var(--muted);
		color: white;
	}

	.role-badge.admin {
		background: var(--brand);
		color: white;
	}

	.actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-icon {
		background: none;
		border: none;
		padding: 0.5rem;
		border-radius: 0.375rem;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.btn-icon:hover {
		transform: scale(1.1);
	}

	.btn-icon.edit {
		color: var(--brand);
	}

	.btn-icon.edit:hover {
		background: rgba(59, 130, 246, 0.1);
	}

	.btn-icon.delete {
		color: var(--danger);
	}

	.btn-icon.delete:hover {
		background: rgba(239, 68, 68, 0.1);
	}

	.btn-icon svg {
		stroke: currentColor;
		fill: none;
		stroke-width: 2;
	}

	/* 모달 스타일 */
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0, 0, 0, 0.5);
		z-index: 9999;
		display: flex;
		justify-content: center;
		align-items: center;
		padding: 20px;
		backdrop-filter: blur(4px);
		width: 100vw;
		height: 100vh;
		min-height: 100vh;
		/* 전역 main 스타일 무시 */
		align-items: center !important;
		justify-content: center !important;
	}

	.modal-content {
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 16px;
		padding: 24px;
		width: 100%;
		max-width: 500px;
		max-height: 80vh;
		overflow-y: auto;
		z-index: 10000;
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 16px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
		animation: modalSlideIn 0.3s ease-out;
		margin: auto;
		align-self: center;
		/* 강제 중앙 정렬 */
		transform: translate(-50%, -50%);
		position: absolute;
		top: 50%;
		left: 50%;
	}

	@keyframes modalSlideIn {
		from {
			opacity: 0;
			transform: translate(-50%, -50%) scale(0.95);
		}
		to {
			opacity: 1;
			transform: translate(-50%, -50%) scale(1);
		}
	}

	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1.5rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid var(--line);
	}

	.modal-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
		color: var(--text);
	}

	.btn-close {
		background: none;
		border: none;
		color: var(--muted);
		cursor: pointer;
		padding: 0.5rem;
		border-radius: 0.375rem;
		transition: all 0.2s ease;
	}

	.btn-close:hover {
		background: var(--bg);
		color: var(--text);
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-bottom: 1rem;
	}

	.form-group label {
		font-weight: 600;
		font-size: 0.875rem;
		color: var(--text);
	}

	.form-input {
		border: 1px solid var(--line);
		background: var(--card);
		color: var(--text);
		padding: 0.75rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		transition: all 0.2s ease;
	}

	.form-input:focus {
		outline: none;
		border-color: var(--brand);
		box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
	}

	.form-input[readonly] {
		background: var(--bg);
		color: var(--muted);
		cursor: not-allowed;
	}

	.form-input textarea {
		resize: vertical;
		min-height: 80px;
		font-family: inherit;
	}

	.form-help {
		font-size: 0.75rem;
		color: var(--muted);
		margin-top: 0.25rem;
		display: block;
	}

	.english-date-badge {
		background: rgba(59, 130, 246, 0.1);
		color: var(--brand);
		padding: 0.25rem 0.5rem;
		border-radius: 0.375rem;
		font-size: 0.75rem;
		font-weight: 600;
		display: inline-flex;
		align-items: center;
		gap: 0.25rem;
	}

	.no-english-date {
		color: var(--muted);
		font-size: 0.875rem;
	}

	.form-actions {
		display: flex;
		gap: 0.75rem;
		justify-content: flex-end;
		margin-top: 1.5rem;
	}

	/* 반응형 디자인 */
	@media (max-width: 768px) {
		.dashboard-grid {
			grid-template-columns: 1fr;
			gap: 1rem;
		}

		.dashboard-header {
			flex-direction: column;
			gap: 1rem;
			align-items: flex-start;
		}

		.header-right {
			width: 100%;
			justify-content: space-between;
		}

		.section-header {
			flex-direction: column;
			gap: 1rem;
			align-items: flex-start;
		}

		.summary-cards {
			grid-template-columns: repeat(2, 1fr);
			padding: 1rem;
		}

		.summary-card {
			padding: 1rem;
		}

		.card-number {
			font-size: 1.5rem;
		}

		.students-table,
		.classes-table {
			font-size: 0.75rem;
		}

			th,
	td {
		padding: 0.75rem 0.5rem;
		white-space: nowrap;
	}

		.actions {
			flex-direction: column;
			gap: 0.25rem;
		}
	}

	@media (max-width: 480px) {
		.summary-cards {
			grid-template-columns: 1fr;
		}

		.summary-header {
			flex-direction: column;
			gap: 1rem;
			align-items: flex-start;
		}
	}

	/* 클릭 가능한 행 스타일 */
	.clickable-row {
		cursor: pointer;
		transition: background-color 0.2s ease;
	}

	.clickable-row:hover {
		background-color: var(--bg);
	}

	/* 시험 기간 모달 스타일 */
	.exam-periods-modal {
		max-width: 800px;
		max-height: 90vh;
	}

	.modal-header-actions {
		display: flex;
		align-items: center;
		gap: 1rem;
	}

	.exam-periods-content {
		max-height: 60vh;
		overflow-y: auto;
	}

	.class-selection {
		padding: 1rem 0;
	}

	.class-selection h4 {
		margin: 0 0 1rem 0;
		color: var(--text);
		font-size: 1.125rem;
		font-weight: 600;
	}

	.class-list {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
		gap: 1rem;
	}

	.class-item {
		background: var(--card);
		border: 1px solid var(--line);
		border-radius: 0.5rem;
		padding: 1rem;
		cursor: pointer;
		transition: all 0.2s ease;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		text-align: left;
	}

	.class-item:hover {
		border-color: var(--brand);
		box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
		transform: translateY(-2px);
	}

	.class-item .class-name {
		font-weight: 600;
		color: var(--text);
		font-size: 1rem;
	}

	.class-item .class-id {
		color: var(--muted);
		font-size: 0.875rem;
	}

	.exam-periods-list {
		padding: 1rem 0;
	}

	.exam-periods-table {
		overflow-x: auto;
	}

	.exam-periods-table table {
		width: 100%;
		border-collapse: collapse;
	}

	.exam-periods-table th,
	.exam-periods-table td {
		padding: 1rem;
		text-align: left;
		border-bottom: 1px solid var(--line);
		white-space: nowrap;
	}

	.exam-periods-table th {
		background: var(--bg);
		font-weight: 600;
		color: var(--text);
		font-size: 0.875rem;
	}

	.exam-periods-table td {
		color: var(--text);
		font-size: 0.875rem;
	}

	.class-link {
		background: none;
		border: none;
		color: var(--brand);
		text-decoration: underline;
		cursor: pointer;
		font-size: 0.875rem;
		padding: 0;
		transition: color 0.2s ease;
	}

	.class-link:hover {
		color: var(--brand-600);
	}

	/* 반응형 디자인 - 시험 기간 모달 */
	@media (max-width: 768px) {
		.exam-periods-modal {
			max-width: 95vw;
			margin: 1rem;
		}

		.modal-header-actions {
			flex-direction: column;
			align-items: stretch;
			gap: 0.5rem;
		}

		.class-list {
			grid-template-columns: 1fr;
		}

		.exam-periods-table th,
		.exam-periods-table td {
			padding: 0.75rem 0.5rem;
			font-size: 0.75rem;
		}
	}
</style>
