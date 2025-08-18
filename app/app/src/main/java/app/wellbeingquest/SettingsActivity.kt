package app.wellbeingquest

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawing
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.Scaffold
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import app.wellbeingquest.ui.theme.BottomBar
import app.wellbeingquest.ui.theme.NavigationButton

class SettingsActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            Scaffold(
                modifier = Modifier.fillMaxSize(),
                contentWindowInsets = WindowInsets.safeDrawing,
                bottomBar = {
                    BottomBar(
                        alignment = Alignment.Start,
                        modifier = Modifier) {
                        NavigationButton(
                            imageVector = Icons.AutoMirrored.Filled.ArrowBack,
                            contentDescription = "Navigate back to Week",
                            onClick = {
                                var intent = Intent(this@SettingsActivity, WeekActivity::class.java)
                                startActivity(intent)
                            },
                        )
                    }
                }
            ) { innerPadding ->
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(innerPadding)
                ) {
                    Greeting(
                        text = "settings",
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }
}