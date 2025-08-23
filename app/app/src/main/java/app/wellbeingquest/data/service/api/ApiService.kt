package app.wellbeingquest.data.service.api

import app.wellbeingquest.data.service.dto.ActivityDto
import app.wellbeingquest.data.service.dto.SuggestionsDto
import app.wellbeingquest.data.service.dto.WeekDto
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path

interface ApiService {
    @POST("activities/v1")
    suspend fun postActivity(@Body activity: ActivityDto): Response<ActivityDto>

    @GET(value="activities/v1/weeks/{weekId}")
    suspend fun getWeek(@Path("weekId") weekId: String): Response<WeekDto>

    @GET(value="activities/v1/suggestions")
    suspend fun getSuggestions(): Response<SuggestionsDto>
}
