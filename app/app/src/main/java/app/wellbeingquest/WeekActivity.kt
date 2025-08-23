package app.wellbeingquest

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawing
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.wrapContentWidth
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowLeft
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowRight
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material.icons.filled.CloudUpload
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.lifecycleScope
import app.wellbeingquest.data.local.database.DatabaseProvider
import app.wellbeingquest.data.model.Activity
import app.wellbeingquest.data.model.Feeling
import app.wellbeingquest.data.service.api.scheduleUploadWorker
import app.wellbeingquest.ui.theme.BottomBar
import app.wellbeingquest.ui.theme.GroupLabel
import app.wellbeingquest.ui.theme.GroupText
import app.wellbeingquest.ui.theme.NavigationButton
import app.wellbeingquest.ui.theme.TopBar
import kotlinx.coroutines.launch
import java.time.LocalDate
import java.time.format.DateTimeFormatter

class WeekActivity : ComponentActivity() {
    var weekNameFormatter: DateTimeFormatter = DateTimeFormatter.ofPattern("yyyy-MM-dd")
    var dayNameFormatter: DateTimeFormatter = DateTimeFormatter.ofPattern("EEE MMM d")

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        lifecycleScope.launch {
            // schedule upload of entries to api service
            scheduleUploadWorker(this@WeekActivity)

            val appDatabase = DatabaseProvider.getInstance(this@WeekActivity)
            val viewModelFactory = WeekViewModelFactory(appDatabase)
            val weekViewModel: WeekViewModel = ViewModelProvider(this@WeekActivity, viewModelFactory)[WeekViewModel::class.java]

            setContent {
                val selectedWeekStart = weekViewModel.selectedWeekStart.collectAsState()
                val hasNextWeek = weekViewModel.hasNextWeek.collectAsState()
                val feelingsInWeek = weekViewModel.feelingsInWeek.collectAsState()

                Scaffold(
                    modifier = Modifier.fillMaxSize(),
                    contentWindowInsets = WindowInsets.safeDrawing,
                    bottomBar = {
                        BottomBar(
                            alignment = Alignment.End,
                            modifier = Modifier
                        ) {
                            NavigationButton(
                                imageVector = Icons.Default.Settings,
                                contentDescription = "Navigate to settings",
                                onClick = {
                                    var intent =
                                        Intent(this@WeekActivity, SettingsActivity::class.java)
                                    startActivity(intent)
                                },
                            )
                            NavigationButton(
                                imageVector = Icons.Default.AddCircle,
                                contentDescription = "Add an activity",
                                onClick = {
                                    var intent = Intent(this@WeekActivity, AddActivity::class.java)
                                    startActivity(intent)
                                },
                                enabled = !hasNextWeek.value,
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
                            arrangement = Arrangement.SpaceBetween,
                            modifier = Modifier
                        ) {

                            NavigationButton(
                                imageVector = Icons.AutoMirrored.Filled.KeyboardArrowLeft,
                                contentDescription = "View the previous week",
                                onClick = {
                                    weekViewModel.previousWeek()
                                },
                            )
                            GroupText(
                                text = getWeekDisplay(selectedWeekStart.value),
                                modifier = Modifier
                                    .wrapContentWidth(Alignment.CenterHorizontally)
                                    .align(Alignment.CenterVertically)
                            )
                            NavigationButton(
                                imageVector = Icons.AutoMirrored.Filled.KeyboardArrowRight,
                                contentDescription = "View the next week",
                                onClick = {
                                    weekViewModel.nextWeek()
                                },
                                enabled = hasNextWeek.value,
                            )

                        }

                        ScrollableContent(feelingsInWeek.value)
                    }
                }
            }
        }
    }

    fun getWeekName(selectedWeek: LocalDate): String {
        return weekNameFormatter.format(selectedWeek)
    }

    fun getWeekDisplay(selectedWeek: LocalDate): String {
        val start = dayNameFormatter.format(selectedWeek)
        val end = dayNameFormatter.format(selectedWeek.plusDays(6))
        return "$start - $end"
    }

    @Composable
    fun ScrollableContent(feelingsInWeek: List<Feeling>) {
        val scrollState = rememberScrollState()

        Column(
            modifier = Modifier
                .verticalScroll(scrollState)
                .padding(8.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp, Alignment.Top)
        ) {
            GroupLabel(
                text = "My feelings and activities",
                modifier = Modifier
                    .fillMaxWidth()
                    .wrapContentWidth(Alignment.CenterHorizontally)
            )
            Spacer(modifier = Modifier.height(4.dp))

            for (feeling in feelingsInWeek) {
                FeelingLabel(feeling.name)
                for (activity in feeling.activities) {
                    ActivityItem(activity)
                }
                Spacer(modifier = Modifier.height(16.dp))
            }
        }
    }

    @Composable
    fun FeelingLabel(feeling: String) {
        Text(
            text = feeling,
            color = Color.White,
            fontSize = 18.sp,
            fontWeight = FontWeight.Bold,
            modifier = Modifier
                .background(
                    color = Color(0xFFE59500), // Gamboge Orange
                    shape = RoundedCornerShape(12.dp)
                )
                .padding(horizontal = 12.dp, vertical = 6.dp)
        )
    }

    @Composable
    fun ActivityItem(activity: Activity) {
        val textColor = if (activity.completed) Color(0xFF02040F) else Color(0xFF6E6E6E)
        val backgroundColor = if (activity.completed) Color(0xFFE5DADA) else Color(0xFFDADADA)
        val fontStyle = if (activity.completed) FontStyle.Normal else FontStyle.Italic
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(start = 16.dp, end = 16.dp)
                .background(
                    color = backgroundColor,
                    shape = RoundedCornerShape(12.dp)
                )
                .padding(horizontal = 12.dp, vertical = 6.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = activity.name,
                color = textColor,
                fontStyle = fontStyle,
                fontWeight = FontWeight.Normal,
                modifier = Modifier.weight(1f)
            )

            if (!activity.completed) {
                Icon(
                    imageVector = Icons.Default.Edit,
                    contentDescription = "Add activity",
                    modifier = Modifier
                        .padding(start = 8.dp)
                        .size(16.dp),
                    tint = textColor
                )
            } else if (!activity.synced) {
                Icon(
                    imageVector = Icons.Filled.CloudUpload,
                    contentDescription = "Sync activity",
                    modifier = Modifier
                        .padding(start = 8.dp)
                        .size(16.dp),
                    tint = textColor
                )
            }
        }
    }
}
