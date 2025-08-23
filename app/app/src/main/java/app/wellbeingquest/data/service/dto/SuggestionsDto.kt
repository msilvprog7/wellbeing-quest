package app.wellbeingquest.data.service.dto

data class SuggestionsDto(
    val activities: List<ActivityDto>,
    val feelings: List<FeelingDto>
)
