package app.wellbeingquest.data.service.api

import app.wellbeingquest.data.service.dto.ActivityDto
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

// todo: add other apis
interface ApiService {
    @POST("activities/v1")
    suspend fun postActivity(@Body activity: ActivityDto): Response<ActivityDto>
}
