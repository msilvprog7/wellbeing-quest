import java.util.Properties

// Setup local.properties
val localProperties = Properties()
val localFile = rootProject.file("local.properties")
if (!localFile.exists()) {
    throw GradleException("Missing local.properties file at project root.")
}

localFile.inputStream().use { localProperties.load(it) }
val googleClientId = localProperties.getProperty("GOOGLE_CLIENT_ID")
    ?: throw GradleException("Missing GOOGLE_CLIENT_ID in local.properties.")
val serviceUrl = localProperties.getProperty("SERVICE_URL")
    ?: throw GradleException("Missing SERVICE_URL in local.properties.")

plugins {
    alias(libs.plugins.android.application)
    alias(libs.plugins.kotlin.android)
    alias(libs.plugins.kotlin.compose)
    id("com.google.devtools.ksp")
}

android {
    namespace = "app.wellbeingquest"
    compileSdk = 36

    buildFeatures {
        buildConfig = true
    }

    defaultConfig {
        applicationId = "app.wellbeingquest"
        minSdk = 36
        targetSdk = 36
        versionCode = 1
        versionName = "1.0"

        testInstrumentationRunner = "androidx.test.runner.AndroidJUnitRunner"

        buildConfigField("String", "GOOGLE_CLIENT_ID", "\"$googleClientId\"")
        buildConfigField("String", "SERVICE_URL", "\"$serviceUrl\"")
    }

    buildTypes {
        release {
            isMinifyEnabled = false
            proguardFiles(
                getDefaultProguardFile("proguard-android-optimize.txt"),
                "proguard-rules.pro"
            )
        }
    }
    compileOptions {
        sourceCompatibility = JavaVersion.VERSION_11
        targetCompatibility = JavaVersion.VERSION_11
    }
    kotlinOptions {
        jvmTarget = "11"
    }
    buildFeatures {
        compose = true
    }
}

dependencies {
    implementation(libs.androidx.core.ktx)
    implementation(libs.androidx.lifecycle.runtime.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.compose)
    implementation(libs.androidx.activity.compose)
    implementation(platform(libs.androidx.compose.bom))
    implementation(libs.androidx.ui)
    implementation(libs.androidx.ui.graphics)
    implementation(libs.androidx.ui.tooling.preview)
    implementation(libs.androidx.material3)
    implementation(libs.androidx.compose.material.iconsCore)
    implementation(libs.androidx.compose.material.iconsExtended)
    implementation(libs.androidx.appcompat)
    implementation(libs.material)
    implementation(libs.androidx.activity)
    implementation(libs.androidx.constraintlayout)
    implementation(libs.material3)
    implementation(libs.androidx.foundation.layout)
    implementation(libs.androidx.room.runtime)
    implementation(libs.androidx.room.ktx)
    implementation(libs.androidx.lifecycle.viewmodel.ktx)
    implementation(libs.androidx.work.runtime.ktx)
    implementation(libs.retrofit2)
    implementation(libs.retrofit2.converter.gson)
    implementation(libs.okhttp3)
    implementation(libs.okhttp3.logging.interceptor)
    implementation(libs.googleid)
    implementation("androidx.credentials:credentials:1.6.0-beta01")
    implementation("androidx.credentials:credentials-play-services-auth:1.6.0-beta01")
    implementation("androidx.datastore:datastore-preferences:1.1.7")
    ksp(libs.androidx.room.compiler)
    testImplementation(libs.junit)
    androidTestImplementation(libs.androidx.junit)
    androidTestImplementation(libs.androidx.espresso.core)
    androidTestImplementation(platform(libs.androidx.compose.bom))
    androidTestImplementation(libs.androidx.ui.test.junit4)
    debugImplementation(libs.androidx.ui.tooling)
    debugImplementation(libs.androidx.ui.test.manifest)
}