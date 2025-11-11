package app.wellbeingquest.data.service.api

import android.content.Intent
import android.os.Handler
import android.os.Looper
import app.wellbeingquest.AuthInterceptorClient
import app.wellbeingquest.AuthInterceptorQualifier
import app.wellbeingquest.LoggingInterceptorQualifier
import app.wellbeingquest.MainActivity
import app.wellbeingquest.NoAuthApiService
import app.wellbeingquest.NoAuthInterceptorClient
import app.wellbeingquest.SessionManager
import app.wellbeingquest.UnauthorizedInterceptorQualifier
import app.wellbeingquest.WellbeingQuestApplication
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import java.util.concurrent.TimeUnit
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.Retrofit
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object RetrofitModule {

    // Default will use auth
    @Provides
    @Singleton
    fun provideApiService(
        sessionManager: SessionManager,
        @AuthInterceptorClient okHttpClient: OkHttpClient
    ): ApiService {
        return Retrofit.Builder()
            .baseUrl(sessionManager.getBaseUrl())
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(ApiService::class.java)
    }

    // Use the @NoAuthApiService to create a Retrofit instance without authentication
    @Provides
    @Singleton
    @NoAuthApiService
    fun provideNoAuthApiService(
        sessionManager: SessionManager,
        @NoAuthInterceptorClient okHttpClient: OkHttpClient
    ): ApiService {
        return Retrofit.Builder()
            .baseUrl(sessionManager.getBaseUrl())
            .client(okHttpClient)
            .addConverterFactory(GsonConverterFactory.create())
            .build()
            .create(ApiService::class.java)
    }

    @Provides
    @Singleton
    @NoAuthInterceptorClient
    fun provideNoAuthOkHttpClient(
        @LoggingInterceptorQualifier loggingInterceptor: HttpLoggingInterceptor
    ): OkHttpClient {
        return OkHttpClient.Builder()
            .addInterceptor(loggingInterceptor)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .writeTimeout(30, TimeUnit.SECONDS)
            .build()
    }

    @Provides
    @Singleton
    @AuthInterceptorClient
    fun provideAuthOkHttpClient(
        @AuthInterceptorQualifier authInterceptor: Interceptor,
        @LoggingInterceptorQualifier loggingInterceptor: HttpLoggingInterceptor,
        @UnauthorizedInterceptorQualifier unauthorizedInterceptor: Interceptor
    ): OkHttpClient {
        return OkHttpClient.Builder()
            .addInterceptor(authInterceptor)
            .addInterceptor(loggingInterceptor)
            .addInterceptor(unauthorizedInterceptor)
            .connectTimeout(30, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .writeTimeout(30, TimeUnit.SECONDS)
            .build()
    }

    @Provides
    @Singleton
    @LoggingInterceptorQualifier
    fun provideLoggingInterceptor(): HttpLoggingInterceptor {
        return HttpLoggingInterceptor().apply {
            level = HttpLoggingInterceptor.Level.BODY
        }
    }

    @Provides
    @Singleton
    @AuthInterceptorQualifier
    fun provideAuthInterceptor(sessionManager: SessionManager): Interceptor {
        return Interceptor { chain ->
            val requestBuilder = chain.request().newBuilder()
            val token = sessionManager.getCachedSessionToken()

            if (!token.isNullOrBlank()) {
                requestBuilder.addHeader("Authorization", "Bearer $token")
            }

            chain.proceed(requestBuilder.build())
        }
    }

    @Provides
    @Singleton
    @UnauthorizedInterceptorQualifier
    fun provideUnauthorizedInterceptor(): Interceptor {
        return Interceptor { chain ->
            val response = chain.proceed(chain.request())
            if (response.code == 401) {
                Handler(Looper.getMainLooper()).post {
                    val context = WellbeingQuestApplication.instance.applicationContext
                    val intent = Intent(context, MainActivity::class.java).apply {
                        addFlags(Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK)
                    }
                    context.startActivity(intent)
                }
            }
            response
        }
    }
}