package service

import (
	"math"
	"sort"

	"go.uber.org/zap"

	"mqfm-backend/internal/domain/entity"
	"mqfm-backend/internal/domain/port"
	"mqfm-backend/internal/shared/logger"
)

const (
	weightPlay = 1.0
	weightLike = 2.0
)

type recommendationService struct {
	audioRepo      port.AudioRepository
	scoreRepo      port.AudioScoreRepository
	historyRepo    port.HistoryRepository
	likeRepo       port.LikeRepository
	locationRepo   port.UserLocationRepository
}

func NewRecommendationService(
	audioRepo port.AudioRepository,
	scoreRepo port.AudioScoreRepository,
	historyRepo port.HistoryRepository,
	likeRepo port.LikeRepository,
	locationRepo port.UserLocationRepository,
) port.RecommendationService {
	return &recommendationService{
		audioRepo:    audioRepo,
		scoreRepo:    scoreRepo,
		historyRepo:  historyRepo,
		likeRepo:     likeRepo,
		locationRepo: locationRepo,
	}
}

func (s *recommendationService) GetPopular(limit int) ([]entity.Audio, error) {
	scores, err := s.scoreRepo.FindTopByScore(limit)
	if err != nil {
		return nil, err
	}
	return extractAudiosFromScores(scores), nil
}

func (s *recommendationService) GetMostListened(limit int) ([]entity.Audio, error) {
	scores, err := s.scoreRepo.FindTopByScore(limit * 2)
	if err != nil {
		return nil, err
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalPlays > scores[j].TotalPlays
	})

	audios := extractAudiosFromScores(scores)
	if len(audios) > limit {
		audios = audios[:limit]
	}
	return audios, nil
}

func (s *recommendationService) GetByArtist(artist string, limit int) ([]entity.Audio, error) {
	return s.audioRepo.FindByArtist(artist, limit)
}

func (s *recommendationService) GetSimilar(audioID uint, limit int) ([]entity.Audio, error) {
	targetAudio, err := s.audioRepo.FindByID(audioID)
	if err != nil {
		return nil, err
	}

	candidates, err := s.audioRepo.FindAllActive()
	if err != nil {
		return nil, err
	}

	audioIDs := make([]uint, 0, len(candidates))
	for _, c := range candidates {
		if c.ID != audioID {
			audioIDs = append(audioIDs, c.ID)
		}
	}

	scoreMap := make(map[uint]*entity.AudioScore)
	if len(audioIDs) > 0 {
		scores, _ := s.scoreRepo.FindByAudioIDs(audioIDs)
		for i := range scores {
			scoreMap[scores[i].AudioID] = &scores[i]
		}
	}

	targetVector := buildFeatureVector(targetAudio, scoreMap[targetAudio.ID])

	type scored struct {
		audio      entity.Audio
		similarity float64
	}

	var scoredList []scored
	for _, c := range candidates {
		if c.ID == audioID {
			continue
		}
		vec := buildFeatureVector(&c, scoreMap[c.ID])
		sim := cosineSimilarity(targetVector, vec)
		scoredList = append(scoredList, scored{audio: c, similarity: sim})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].similarity > scoredList[j].similarity
	})

	var result []entity.Audio
	for i, s := range scoredList {
		if i >= limit {
			break
		}
		result = append(result, s.audio)
	}

	return result, nil
}

func (s *recommendationService) GetQuickPick(userID uint, limit int) ([]entity.Audio, error) {
	histories, err := s.historyRepo.FindFrequentByUser(userID, 2, limit)
	if err != nil {
		return nil, err
	}

	var audios []entity.Audio
	for _, h := range histories {
		if h.Audio != nil {
			audios = append(audios, *h.Audio)
		}
	}

	return audios, nil
}

func (s *recommendationService) GetOnboarding(limit int) ([]entity.Audio, error) {
	scores, err := s.scoreRepo.FindTopByScore(limit * 3)
	if err != nil || len(scores) == 0 {
		return s.audioRepo.FindAllActive()
	}

	categorySet := make(map[uint]bool)
	var diverseAudios []entity.Audio

	for _, sc := range scores {
		if sc.Audio == nil {
			continue
		}
		if !categorySet[sc.Audio.CategoryID] {
			categorySet[sc.Audio.CategoryID] = true
			diverseAudios = append(diverseAudios, *sc.Audio)
		}
		if len(diverseAudios) >= limit {
			break
		}
	}

	if len(diverseAudios) < limit {
		for _, sc := range scores {
			if sc.Audio == nil {
				continue
			}
			found := false
			for _, a := range diverseAudios {
				if a.ID == sc.Audio.ID {
					found = true
					break
				}
			}
			if !found {
				diverseAudios = append(diverseAudios, *sc.Audio)
			}
			if len(diverseAudios) >= limit {
				break
			}
		}
	}

	return diverseAudios, nil
}

func (s *recommendationService) GetPersonalized(userID uint, limit int) ([]entity.Audio, error) {
	histories, err := s.historyRepo.FindByUser(userID)
	if err != nil || len(histories) == 0 {
		return s.GetOnboarding(limit)
	}

	likes, _ := s.likeRepo.FindByUser(userID)

	listenedIDs := make(map[uint]bool)
	var listenedAudios []*entity.Audio
	for _, h := range histories {
		listenedIDs[h.AudioID] = true
		if h.Audio != nil {
			listenedAudios = append(listenedAudios, h.Audio)
		}
	}
	for _, l := range likes {
		listenedIDs[l.AudioID] = true
		if l.Audio != nil {
			listenedAudios = append(listenedAudios, l.Audio)
		}
	}

	if len(listenedAudios) == 0 {
		return s.GetOnboarding(limit)
	}

	allAudioIDs := make([]uint, 0)
	for _, a := range listenedAudios {
		allAudioIDs = append(allAudioIDs, a.ID)
	}
	scoreMap := make(map[uint]*entity.AudioScore)
	scores, _ := s.scoreRepo.FindByAudioIDs(allAudioIDs)
	for i := range scores {
		scoreMap[scores[i].AudioID] = &scores[i]
	}

	userVector := [4]float64{}
	for _, a := range listenedAudios {
		vec := buildFeatureVector(a, scoreMap[a.ID])
		for i := range userVector {
			userVector[i] += vec[i]
		}
	}
	n := float64(len(listenedAudios))
	for i := range userVector {
		userVector[i] /= n
	}

	candidates, err := s.audioRepo.FindAllActive()
	if err != nil {
		return nil, err
	}

	candIDs := make([]uint, 0, len(candidates))
	for _, c := range candidates {
		if !listenedIDs[c.ID] {
			candIDs = append(candIDs, c.ID)
		}
	}
	candScores, _ := s.scoreRepo.FindByAudioIDs(candIDs)
	candScoreMap := make(map[uint]*entity.AudioScore)
	for i := range candScores {
		candScoreMap[candScores[i].AudioID] = &candScores[i]
	}

	type scored struct {
		audio      entity.Audio
		similarity float64
	}

	var scoredList []scored
	for _, c := range candidates {
		if listenedIDs[c.ID] {
			continue
		}
		vec := buildFeatureVector(&c, candScoreMap[c.ID])
		sim := cosineSimilarity(userVector, vec)
		scoredList = append(scoredList, scored{audio: c, similarity: sim})
	}

	sort.Slice(scoredList, func(i, j int) bool {
		return scoredList[i].similarity > scoredList[j].similarity
	})

	var result []entity.Audio
	for i, s := range scoredList {
		if i >= limit {
			break
		}
		result = append(result, s.audio)
	}

	if len(result) == 0 {
		return s.GetOnboarding(limit)
	}

	return result, nil
}

func (s *recommendationService) RecalculateScores() error {
	playCounts, err := s.historyRepo.AggregatePlayCounts()
	if err != nil {
		logger.Error("recalculate: failed to aggregate play counts", zap.Error(err))
		return err
	}

	likeCounts, err := s.likeRepo.AggregateLikeCounts()
	if err != nil {
		logger.Error("recalculate: failed to aggregate like counts", zap.Error(err))
		return err
	}

	audioIDSet := make(map[uint]bool)
	for id := range playCounts {
		audioIDSet[id] = true
	}
	for id := range likeCounts {
		audioIDSet[id] = true
	}

	var scores []entity.AudioScore
	for audioID := range audioIDSet {
		plays := playCounts[audioID]
		likes := likeCounts[audioID]
		weight := float64(plays)*weightPlay + float64(likes)*weightLike

		scores = append(scores, entity.AudioScore{
			AudioID:     audioID,
			TotalPlays:  plays,
			TotalLikes:  likes,
			WeightScore: weight,
		})
	}

	if err := s.scoreRepo.BulkUpsert(scores); err != nil {
		logger.Error("recalculate: bulk upsert failed", zap.Error(err))
		return err
	}

	logger.Info("recalculate: scores updated", zap.Int("count", len(scores)))
	return nil
}

func buildFeatureVector(audio *entity.Audio, score *entity.AudioScore) [4]float64 {
	var plays, likes float64
	if score != nil {
		plays = float64(score.TotalPlays)
		likes = float64(score.TotalLikes)
	}

	normalizedPlays := math.Log1p(plays)
	normalizedLikes := math.Log1p(likes)
	categoryFeature := float64(audio.CategoryID % 100)
	artistHash := float64(simpleHash(audio.Artist) % 1000)

	return [4]float64{normalizedPlays, normalizedLikes, categoryFeature, artistHash}
}

func cosineSimilarity(a, b [4]float64) float64 {
	var dot, normA, normB float64
	for i := 0; i < 4; i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func simpleHash(s string) uint64 {
	var h uint64 = 5381
	for _, c := range s {
		h = h*33 + uint64(c)
	}
	return h
}

func extractAudiosFromScores(scores []entity.AudioScore) []entity.Audio {
	var audios []entity.Audio
	for _, s := range scores {
		if s.Audio != nil {
			audios = append(audios, *s.Audio)
		}
	}
	return audios
}

func (s *recommendationService) GetLocationBased(userID uint, limit int) ([]entity.Audio, error) {
	loc, err := s.locationRepo.FindByUser(userID)
	if err != nil || loc.City == "" {
		return s.GetPopular(limit)
	}

	audios, err := s.audioRepo.Search(loc.City)
	if err != nil || len(audios) == 0 {
		return s.GetPopular(limit)
	}

	if len(audios) > limit {
		audios = audios[:limit]
	}
	return audios, nil
}

func (s *recommendationService) GetTimeBasedPersonalized(userID uint, hour int, limit int) ([]entity.Audio, error) {
	switch {
	case hour >= 5 && hour < 12:
		return s.GetPopular(limit)
	case hour >= 12 && hour < 17:
		return s.GetMostListened(limit)
	case hour >= 17 && hour < 21:
		return s.GetPersonalized(userID, limit)
	default:
		return s.GetQuickPick(userID, limit)
	}
}
