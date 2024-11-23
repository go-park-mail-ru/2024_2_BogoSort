package survey

import (
	pb "github.com/go-park-mail-ru/2024_2_BogoSort/internal/delivery/grpc/survey/proto"
	"github.com/go-park-mail-ru/2024_2_BogoSort/internal/entity/dto"
	"github.com/pkg/errors"
)

func ConvertDBPageTypeToEnum(dbPageType string) (pb.PageType, error) {
	switch dbPageType {
	case "mainPage":
		return pb.PageType_PAGE_TYPE_MAIN, nil
	case "advertPage":
		return pb.PageType_PAGE_TYPE_ADVERT, nil
	case "advertCreatePage":
		return pb.PageType_PAGE_TYPE_ADVERT_CREATE, nil
	case "cartPage":
		return pb.PageType_PAGE_TYPE_CART, nil
	case "categoryPage":
		return pb.PageType_PAGE_TYPE_CATEGORY, nil
	case "advertEditPage":
		return pb.PageType_PAGE_TYPE_ADVERT_EDIT, nil
	case "userPage":
		return pb.PageType_PAGE_TYPE_USER, nil
	case "sellerPage":
		return pb.PageType_PAGE_TYPE_SELLER, nil
	case "searchPage":
		return pb.PageType_PAGE_TYPE_SEARCH, nil
	default:
		return pb.PageType_PAGE_TYPE_UNSPECIFIED, errors.New("unknown page type")
	}
}

func ConvertEnumToDBPageType(pageType pb.PageType) string {
	switch pageType {
	case pb.PageType_PAGE_TYPE_MAIN:
		return "mainPage"
	case pb.PageType_PAGE_TYPE_ADVERT:
		return "advertPage"
	case pb.PageType_PAGE_TYPE_ADVERT_CREATE:
		return "advertCreatePage"
	case pb.PageType_PAGE_TYPE_CART:
		return "cartPage"
	case pb.PageType_PAGE_TYPE_CATEGORY:
		return "categoryPage"
	case pb.PageType_PAGE_TYPE_ADVERT_EDIT:
		return "advertEditPage"
	case pb.PageType_PAGE_TYPE_USER:
		return "userPage"
	case pb.PageType_PAGE_TYPE_SELLER:
		return "sellerPage"
	case pb.PageType_PAGE_TYPE_SEARCH:
		return "searchPage"
	default:
		return "unknown"
	}
}

func ConvertDBAnswerStatsToProto(answerStats []dto.AnswerStats) []*pb.AnswerStats {
	protoAnswerStats := []*pb.AnswerStats{}
	for _, answerStats := range answerStats {
		protoAnswerStats = append(protoAnswerStats, &pb.AnswerStats{
			Value: int32(answerStats.Value),
			Count: int32(answerStats.Count),
		})
	}
	return protoAnswerStats
}

func ConvertDBQuestionStatsToProto(questionStats []dto.QuestionStats) []*pb.QuestionStats {
	protoQuestionStats := []*pb.QuestionStats{}
	for _, questionStats := range questionStats {
		protoQuestionStats = append(protoQuestionStats, &pb.QuestionStats{
			AnswerStats: ConvertDBAnswerStatsToProto(questionStats.AnswerStats),
			AvgValue:    int32(questionStats.AvgValue),
			Title:       questionStats.Title,
		})
	}
	return protoQuestionStats
}

func ConvertDBStatsToProto(stats *dto.GetStatsResponse) []*pb.PageStats {
	protoStats := []*pb.PageStats{}
	for _, pageStats := range stats.PageStats {
		protoStats = append(protoStats, &pb.PageStats{
			Page: pageStats.Page,
			QuestionStats: ConvertDBQuestionStatsToProto(pageStats.QuestionStats),
		})
	}
	return protoStats
}

func ConvertProtoStatsToDB(protoStats *pb.GetStatsResponse) *dto.GetStatsResponse {
	stats := &dto.GetStatsResponse{}
	for _, pageStats := range protoStats.PageStats {
		stats.PageStats = append(stats.PageStats, dto.PageStats{
			Page: pageStats.Page,
		})
	}
	return stats
}
