package app.wellbeingquest.data.service.dto

data class ActivityDto(
    val name: String,
    val feelings: List<String>? = null,
    val week: String? = null,
    val created: String? = null
)
