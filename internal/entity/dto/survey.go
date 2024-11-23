package dto

type GetQuestionsRequest struct {
	Page string `json:"page"`
}

type PostAnswersRequest struct {
	Page    string   `json:"page"`
	Answers []Answer `json:"answers"`
}
