package lms

import (
	"github.com/c-learn/pkg/errs"
	"github.com/c-learn/pkg/response"
	"github.com/c-learn/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo *Repository
}

func NewAdminHandler(repo *Repository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

// ========== Modules ==========

type ModuleRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

func (h *AdminHandler) ListModules(c *gin.Context) {
	modules, err := h.repo.ListModules(c.Request.Context())
	if err != nil {
		response.Error(c, errs.NewInternal("list modules failed"))
		return
	}
	response.OK(c, modules)
}

type AdminTestCaseResponse struct {
	ID        string `json:"id"`
	Input     string `json:"input"`
	Expected  string `json:"expected"`
	IsSample  bool   `json:"is_sample"`
	SortOrder int    `json:"sort_order"`
}

type AdminTaskResponse struct {
	ID          string                  `json:"id"`
	LessonID    string                  `json:"lesson_id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Difficulty  string                  `json:"difficulty"`
	SortOrder   int                     `json:"sort_order"`
	TestCases   []AdminTestCaseResponse `json:"test_cases"`
}

type AdminLessonResponse struct {
	ID        string              `json:"id"`
	ModuleID  string              `json:"module_id"`
	Title     string              `json:"title"`
	TheoryMD  string              `json:"theory_md"`
	SortOrder int                 `json:"sort_order"`
	Tasks     []AdminTaskResponse `json:"tasks"`
}

type AdminModuleFullResponse struct {
	ID          string                `json:"id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	SortOrder   int                   `json:"sort_order"`
	Lessons     []AdminLessonResponse `json:"lessons"`
}

func (h *AdminHandler) GetModuleFull(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()

	modules, err := h.repo.ListModules(ctx)
	if err != nil {
		response.Error(c, errs.NewInternal("list modules failed"))
		return
	}
	var found *struct {
		ID, Title, Description string
		SortOrder              int
	}
	for _, m := range modules {
		if m.ID == id {
			found = &struct {
				ID, Title, Description string
				SortOrder              int
			}{m.ID, m.Title, m.Description, m.SortOrder}
			break
		}
	}
	if found == nil {
		response.Error(c, errs.NewNotFound("module not found"))
		return
	}

	lessons, err := h.repo.ListLessonsByModule(ctx, id)
	if err != nil {
		response.Error(c, errs.NewInternal("list lessons failed"))
		return
	}

	lessonResps := make([]AdminLessonResponse, 0, len(lessons))
	for _, l := range lessons {
		tasks, err := h.repo.ListTasksByLesson(ctx, l.ID)
		if err != nil {
			response.Error(c, errs.NewInternal("list tasks failed"))
			return
		}
		taskResps := make([]AdminTaskResponse, 0, len(tasks))
		for _, t := range tasks {
			tcs, err := h.repo.GetTestCasesByTask(ctx, t.ID)
			if err != nil {
				response.Error(c, errs.NewInternal("get test cases failed"))
				return
			}
			tcResps := make([]AdminTestCaseResponse, 0, len(tcs))
			for _, tc := range tcs {
				tcResps = append(tcResps, AdminTestCaseResponse{
					ID:        tc.ID,
					Input:     tc.Input,
					Expected:  tc.Expected,
					IsSample:  tc.IsSample,
					SortOrder: tc.SortOrder,
				})
			}
			taskResps = append(taskResps, AdminTaskResponse{
				ID:          t.ID,
				LessonID:    t.LessonID,
				Title:       t.Title,
				Description: t.Description,
				Difficulty:  t.Difficulty,
				SortOrder:   t.SortOrder,
				TestCases:   tcResps,
			})
		}
		lessonResps = append(lessonResps, AdminLessonResponse{
			ID:        l.ID,
			ModuleID:  l.ModuleID,
			Title:     l.Title,
			TheoryMD:  l.TheoryMD,
			SortOrder: l.SortOrder,
			Tasks:     taskResps,
		})
	}

	response.OK(c, AdminModuleFullResponse{
		ID:          found.ID,
		Title:       found.Title,
		Description: found.Description,
		SortOrder:   found.SortOrder,
		Lessons:     lessonResps,
	})
}

func (h *AdminHandler) CreateModule(c *gin.Context) {
	var req ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	if err := validator.ValidateRequired(req.Title, "title"); err != nil {
		response.Error(c, errs.NewValidation(err.Error(), nil))
		return
	}
	m, err := h.repo.CreateModule(c.Request.Context(), req.Title, req.Description, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("create module failed"))
		return
	}
	response.Created(c, m)
}

func (h *AdminHandler) UpdateModule(c *gin.Context) {
	id := c.Param("id")
	var req ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	m, err := h.repo.UpdateModule(c.Request.Context(), id, req.Title, req.Description, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("update module failed"))
		return
	}
	if m == nil {
		response.Error(c, errs.NewNotFound("module not found"))
		return
	}
	response.OK(c, m)
}

func (h *AdminHandler) DeleteModule(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteModule(c.Request.Context(), id); err != nil {
		response.Error(c, errs.NewNotFound("module not found"))
		return
	}
	response.NoContent(c)
}

// ========== Lessons ==========

type LessonRequest struct {
	ModuleID  string `json:"module_id"`
	Title     string `json:"title"`
	TheoryMD  string `json:"theory_md"`
	SortOrder int    `json:"sort_order"`
}

func (h *AdminHandler) CreateLesson(c *gin.Context) {
	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	if err := validator.ValidateRequired(req.Title, "title"); err != nil {
		response.Error(c, errs.NewValidation(err.Error(), nil))
		return
	}
	l, err := h.repo.CreateLesson(c.Request.Context(), req.ModuleID, req.Title, req.TheoryMD, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("create lesson failed"))
		return
	}
	response.Created(c, l)
}

func (h *AdminHandler) UpdateLesson(c *gin.Context) {
	id := c.Param("id")
	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	l, err := h.repo.UpdateLesson(c.Request.Context(), id, req.Title, req.TheoryMD, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("update lesson failed"))
		return
	}
	if l == nil {
		response.Error(c, errs.NewNotFound("lesson not found"))
		return
	}
	response.OK(c, l)
}

func (h *AdminHandler) DeleteLesson(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteLesson(c.Request.Context(), id); err != nil {
		response.Error(c, errs.NewNotFound("lesson not found"))
		return
	}
	response.NoContent(c)
}

// ========== Tasks ==========

type TaskRequest struct {
	LessonID    string            `json:"lesson_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Difficulty  string            `json:"difficulty"`
	SortOrder   int               `json:"sort_order"`
	TestCases   []TestCaseRequest `json:"test_cases"`
}

type TestCaseRequest struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	IsSample bool   `json:"is_sample"`
}

func (h *AdminHandler) CreateTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	if err := validator.ValidateRequired(req.Title, "title"); err != nil {
		response.Error(c, errs.NewValidation(err.Error(), nil))
		return
	}

	task, err := h.repo.CreateTask(c.Request.Context(), req.LessonID, req.Title, req.Description, req.Difficulty, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("create task failed"))
		return
	}

	for i, tc := range req.TestCases {
		if _, err := h.repo.CreateTestCase(c.Request.Context(), task.ID, tc.Input, tc.Expected, i, tc.IsSample); err != nil {
			response.Error(c, errs.NewInternal("create test case failed"))
			return
		}
	}

	response.Created(c, task)
}

func (h *AdminHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}
	task, err := h.repo.UpdateTask(c.Request.Context(), id, req.Title, req.Description, req.Difficulty, req.SortOrder)
	if err != nil {
		response.Error(c, errs.NewInternal("update task failed"))
		return
	}
	if task == nil {
		response.Error(c, errs.NewNotFound("task not found"))
		return
	}
	response.OK(c, task)
}

func (h *AdminHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteTask(c.Request.Context(), id); err != nil {
		response.Error(c, errs.NewNotFound("task not found"))
		return
	}
	response.NoContent(c)
}

type UpdateTestCasesRequest struct {
	TestCases []TestCaseRequest `json:"test_cases"`
}

func (h *AdminHandler) UpdateTestCases(c *gin.Context) {
	taskID := c.Param("id")
	var req UpdateTestCasesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	if err := h.repo.DeleteTestCasesByTask(c.Request.Context(), taskID); err != nil {
		response.Error(c, errs.NewInternal("delete test cases failed"))
		return
	}

	for i, tc := range req.TestCases {
		if _, err := h.repo.CreateTestCase(c.Request.Context(), taskID, tc.Input, tc.Expected, i, tc.IsSample); err != nil {
			response.Error(c, errs.NewInternal("create test case failed"))
			return
		}
	}

	response.OK(c, gin.H{"message": "test cases updated"})
}
