package app.wellbeingquest

import javax.inject.Qualifier

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class AuthInterceptorClient

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class NoAuthInterceptorClient

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class NoAuthApiService

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class AuthInterceptorQualifier

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class UnauthorizedInterceptorQualifier

@Qualifier
@Retention(AnnotationRetention.BINARY)
annotation class LoggingInterceptorQualifier