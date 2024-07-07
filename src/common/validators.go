package common

func IsVideoValid(video *Video) bool {
	if len(video.FileType) == 0 {
		return false
	}
	if video.FileSize == 0 {
		return false
	}
	if video.CreationTimestamp == 0 {
		return false
	}
	if video.LastChangeTimestamp == 0 {
		return false
	}
	return true
}

func IsEpisodeValid(episode *Episode) bool {
	if episode.EpisodeNumber == 0 {
		return false
	}
	if len(episode.Title) == 0 {
		return false
	}
	if len(episode.Description) == 0 {
		return false
	}
	for _, actor := range episode.Actors {
		if len(actor) == 0 {
			return false
		}
	}
	for _, director := range episode.Directors {
		if len(director) == 0 {
			return false
		}
	}
	return IsVideoValid(&episode.Video)
}

func IsSeasonValid(season *Season) bool {
	if season.SeasonNumber == 0 {
		return false
	}
	for _, episode := range season.Episodes {
		if !IsEpisodeValid(&episode) {
			return false
		}
	}
	return true
}

func IsMovieValid(movie *Movie) bool {
	if len(movie.Title) == 0 {
		return false
	}
	if len(movie.Description) == 0 {
		return false
	}
	for _, genre := range movie.Genres {
		if len(genre) == 0 {
			return false
		}
	}
	for _, actor := range movie.Actors {
		if len(actor) == 0 {
			return false
		}
	}
	for _, director := range movie.Directors {
		if len(director) == 0 {
			return false
		}
	}
	return IsVideoValid(&movie.Video)
}

func IsShowValid(show *Show) bool {
	if len(show.Title) == 0 {
		return false
	}
	if len(show.Description) == 0 {
		return false
	}
	for _, genre := range show.Genres {
		if len(genre) == 0 {
			return false
		}
	}
	for _, actor := range show.Actors {
		if len(actor) == 0 {
			return false
		}
	}
	for _, director := range show.Directors {
		if len(director) == 0 {
			return false
		}
	}
	for _, season := range show.Seasons {
		if !IsSeasonValid(&season) {
			return false
		}
	}
	return true
}

func (subscription *Subscription) IsValid() bool {
	if len(subscription.UserUUID) == 0 {
		return false
	}
	if len(subscription.Target) == 0 {
		return false
	}
	if subscription.Type < 0 || subscription.Type > 3 {
		return false
	}

	return true
}

func (editMetadataRequest *EditMetadataRequest) IsValid() bool {
	if len(editMetadataRequest.TargetUUID) == 0 {
		return false
	}
	if len(editMetadataRequest.Title) == 0 {
		return false
	}
	for _, genre := range editMetadataRequest.Genres {
		if len(genre) == 0 {
			return false
		}
	}
	for _, actor := range editMetadataRequest.Actors {
		if len(actor) == 0 {
			return false
		}
	}
	if editMetadataRequest.SeasonNumber != nil &&
		editMetadataRequest.EpisodeNumber != nil &&
		len(editMetadataRequest.Directors) != 0 {
		return false
	}
	if editMetadataRequest.SeasonNumber != nil && editMetadataRequest.EpisodeNumber == nil ||
		editMetadataRequest.SeasonNumber == nil && editMetadataRequest.EpisodeNumber != nil {
		return false
	}

	return true
}
