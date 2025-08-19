package app.wellbeingquest

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawing
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import app.wellbeingquest.ui.theme.AutoCompleteTextField
import app.wellbeingquest.ui.theme.BottomBar
import app.wellbeingquest.ui.theme.FormButton
import app.wellbeingquest.ui.theme.GroupText
import app.wellbeingquest.ui.theme.MultiEntryTextField
import app.wellbeingquest.ui.theme.NavigationButton
import app.wellbeingquest.ui.theme.TopBar

class AddActivity : ComponentActivity() {
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
                                var intent = Intent(this@AddActivity, WeekActivity::class.java)
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
                    TopBar(
                        arrangement = Arrangement.Center,
                        modifier = Modifier
                    ) {
                        GroupText("add an activity", modifier = Modifier)
                    }

                    ScrollableContent()
                }
            }
        }
    }

    @Composable
    fun ScrollableContent() {
        val scrollState = rememberScrollState()
        var activity = remember { mutableStateOf("") }
        var feeling = remember { mutableStateOf("") }
        var feelings by remember { mutableStateOf(listOf<String>()) }
        var activitySuggestions = remember(activity) {
            var suggestions = listOf("Exercise", "Read", "Meditate", "Listen to music")
            suggestions.filter { suggestion ->
                suggestion.contains(activity.value, ignoreCase = true) && suggestion.lowercase() != activity.value
            }.take(5)
        }
        var activitySuggestionsExpanded by remember { mutableStateOf(false) }
        var feelingSuggestions = remember(feeling) {
            var suggestions = listOf("Relaxed", "Focused", "Energetic", "Sad")
            suggestions.filter { suggestion ->
                suggestion.contains(feeling.value, ignoreCase = true) && suggestion.lowercase() != feeling.value
            }.take(5)
        }
        var feelingSuggestionsExpanded by remember { mutableStateOf(false) }

        Column(
            modifier = Modifier
                .verticalScroll(scrollState)
                .padding(8.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp, Alignment.Top)
        ) {
            // activity name
            AutoCompleteTextField(
                text = activity.value,
                onTextChange = { activity.value = it },
                label = { Text("enter name of activity") },
                placeholder = { Text("Activity") },
                suggestions = activitySuggestions,
                onSuggestionsChange = { activitySuggestions = it },
                expanded = activitySuggestionsExpanded,
                onExpandedChange = { activitySuggestionsExpanded = it },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(16.dp)
            )
            Spacer(modifier = Modifier.height(16.dp))

            // feelings
            MultiEntryTextField(
                text = feeling.value,
                onTextChange = { feeling.value = it },
                label = { Text("enter feelings") },
                placeholder = { Text("Relaxed") },
                suggestions = feelingSuggestions,
                onSuggestionsChange = { feelingSuggestions = it },
                expanded = feelingSuggestionsExpanded,
                onExpandedChange = { feelingSuggestionsExpanded = it },
                selected = feelings,
                onSelectedChange = { feelings = it },
                modifier = Modifier
            )
            Spacer(modifier = Modifier.height(16.dp))

            // add button
            FormButton(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "add activity",
                onClick = {
                  var intent = Intent(this@AddActivity, WeekActivity::class.java)
                  startActivity(intent)
                },
                // require name and 1 feeling
                enabled = activity.value.isNotBlank() && feelings.isNotEmpty(),
            )
        }
    }
}