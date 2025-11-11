package app.wellbeingquest

import android.content.Intent
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.BackHandler
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material3.Button
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.credentials.CredentialManager
import androidx.credentials.GetCredentialRequest
import androidx.credentials.exceptions.GetCredentialException
import androidx.lifecycle.lifecycleScope
import app.wellbeingquest.data.service.api.ApiService
import app.wellbeingquest.data.service.dto.IdTokenRequest
import app.wellbeingquest.ui.theme.WellbeingQuestTheme
import com.google.android.libraries.identity.googleid.GetSignInWithGoogleOption
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @Inject
    @NoAuthApiService
    lateinit var api: ApiService

    @Inject
    lateinit var sessionManager: SessionManager

    @Inject
    lateinit var credentialManager: CredentialManager

    private val TAG = "MainActivity"
    private val SERVER_CLIENT_ID = BuildConfig.GOOGLE_CLIENT_ID;


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        enableEdgeToEdge()

        setContent {
            WellbeingQuestTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    Box(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding)
                    ) {
                        Column(
                            modifier = Modifier
                                .fillMaxSize()
                                .align(Alignment.Center),
                            horizontalAlignment = Alignment.CenterHorizontally,
                            verticalArrangement = Arrangement.Center
                        ) {
                            Greeting(
                                text = "wellbeing quest",
                                modifier = Modifier.align(Alignment.CenterHorizontally)
                            )
                        }

                        Button(
                            onClick = { initiateGoogleSignIn() },
                            modifier = Modifier
                                .align(Alignment.BottomCenter)
                                .padding(bottom = 32.dp)
                        ) {
                            Icon(
                                imageVector = Icons.Default.AccountCircle,
                                contentDescription = "Sign in with Google",
                                modifier = Modifier.size(16.dp)
                            )
                            Spacer(modifier = Modifier.width(8.dp))
                            Text(text = "Sign in with Google")
                        }
                    }

                    BackHandler(enabled = true) {
                        // Prevent back navigation
                    }
                }
            }
        }
    }

    private fun initiateGoogleSignIn() {
        lifecycleScope.launch {
            try {
                Log.d(TAG, "Initiating Google sign-in. Server client id: ${SERVER_CLIENT_ID}")
                val googleSignInOption = GetSignInWithGoogleOption.Builder(SERVER_CLIENT_ID)
                    .build()

                val request = GetCredentialRequest.Builder()
                    .addCredentialOption(googleSignInOption)
                    .build()

                val result = credentialManager.getCredential(
                    request = request,
                    context = this@MainActivity
                )

                Log.d(TAG, "Received credential. Type: ${result.credential.type}")
                val idTokenCredential = GoogleIdTokenCredential.createFrom(result.credential.data);

                if (idTokenCredential?.idToken == null) {
                    Log.w("SignIn", "No ID token received")
                    return@launch
                }

                Log.d("SignIn", "ID Token: ${idTokenCredential.idToken}")

                // validate in backend
                val sessionTokenResponse = api.authGoogle(
                    IdTokenRequest(
                        idTokenCredential.idToken!!
                    )
                )

                if (!sessionTokenResponse.isSuccessful || sessionTokenResponse.body() == null || sessionTokenResponse.body()?.sessionToken == null) {
                    Log.d("SignIn", "No Session token received")
                }

                val sessionToken = sessionTokenResponse.body()?.sessionToken!!
                Log.d("SignIn", "Session token: $sessionToken")

                // store token
                sessionManager.saveSessionToken(sessionToken)

                // send to next screen
                var intent = Intent(this@MainActivity, WeekActivity::class.java)
                startActivity(intent)
            } catch (e: GetCredentialException) {
                Log.e("SignIn", "Sign-in failed: ${e.message}", e)
            }
        }
    }
}

@Composable
fun Greeting(text: String, modifier: Modifier = Modifier) {
    Box(
        contentAlignment = Alignment.Center,
        modifier = modifier.fillMaxSize()
    ) {
        Text(
            text = text,
            fontSize = 24.sp
        )
    }
}
