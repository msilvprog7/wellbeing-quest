package app.wellbeingquest.data.service.dto

import java.time.LocalDate

data class ActivityDto(
    val name: String,
    val feelings: List<String>? = null,
    val week: String? = null,
    val created: String? = null
)
