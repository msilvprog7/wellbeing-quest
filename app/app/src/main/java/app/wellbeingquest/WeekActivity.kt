package app.wellbeingquest

import android.os.Bundle
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.ScrollState
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.WindowInsets
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.safeDrawing
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.layout.wrapContentWidth
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowLeft
import androidx.compose.material.icons.automirrored.filled.KeyboardArrowRight
import androidx.compose.material.icons.filled.AddCircle
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Settings
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

class WeekActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            Scaffold(
                modifier = Modifier.fillMaxSize(),
                contentWindowInsets = WindowInsets.safeDrawing,
                bottomBar = {
                    BottomBar()
                }
            ) { innerPadding ->
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(innerPadding)
                ) {
                    TopBar()
                    ScrollableContent()
                }
            }
        }
    }

    @Composable
    fun TopBar() {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(8.dp),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            NavigationButton(
                imageVector = Icons.AutoMirrored.Filled.KeyboardArrowLeft,
                contentDescription = "View the previous week"
            )
            GroupText(
                text = "Sun Aug 10 - Sat Aug 16",
                modifier = Modifier
                    .wrapContentWidth(Alignment.CenterHorizontally)
                    .align(Alignment.CenterVertically))
            NavigationButton(
                imageVector = Icons.AutoMirrored.Filled.KeyboardArrowRight,
                contentDescription = "View the next week"
            )
        }
    }

    @Composable
    fun ScrollableContent() {
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
                    .wrapContentWidth(Alignment.CenterHorizontally))
            Spacer(modifier = Modifier.height(4.dp))

            FeelingLabel("Relaxed feeling")
            ActivityItem("Meditated")
            ActivityItem("Read")
            ActivityItem("Coffee", incomplete = true)
            Spacer(modifier = Modifier.height(16.dp))

            FeelingLabel("Accomplished feeling")
            ActivityItem("Chores")
            ActivityItem("Gardened")
            ActivityItem("Played with cat", incomplete = true)
            ActivityItem("Met up with friends", incomplete = true)
            ActivityItem("Exercised at gym", incomplete = true)
            Spacer(modifier = Modifier.height(16.dp))

            FeelingLabel("Excited feeling")
            ActivityItem("Worked on a new project")
            ActivityItem("Read an interesting article")
            ActivityItem("Found a new snack for my cat")
            ActivityItem("Heard a new song on the radio", incomplete = true)
            ActivityItem("Learned a new fact about space", incomplete = true)
            ActivityItem("Watched a new thriller movie", incomplete = true)
            ActivityItem("Signed up for a ceramics class", incomplete = true)
            ActivityItem("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", incomplete = true)
            Spacer(modifier = Modifier.height(16.dp))
        }
    }

    @Composable
    fun BottomBar() {
        // Bottom Bar
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(8.dp),
            horizontalArrangement = Arrangement.spacedBy(16.dp, Alignment.End),
            verticalAlignment = Alignment.CenterVertically
        ) {
            NavigationButton(
                imageVector = Icons.Default.Settings,
                contentDescription = "Navigate to settings"
            )
            NavigationButton(
                imageVector = Icons.Default.AddCircle,
                contentDescription = "Add an activity"
            )
        }
    }

    @Composable
    fun GroupText(text: String, modifier: Modifier) {
        Text(
            text = "Sun Aug 10 - Sat Aug 16",
            fontSize = 20.sp,
            fontWeight = FontWeight.Bold,
            modifier = modifier
        )
    }

    @Composable
    fun GroupLabel(text: String, modifier: Modifier) {
        Text(
            text = text,
            color = Color.White,
            fontSize = 20.sp,
            fontWeight = FontWeight.Bold,
            modifier = modifier
                .background(
                    color = Color(0xFF002642), // Prussian Blue
                    shape = RoundedCornerShape(12.dp)
                )
                .padding(horizontal = 12.dp, vertical = 6.dp)
        )
    }

    @Composable
    fun NavigationButton(imageVector: ImageVector, contentDescription: String) {
        var context = LocalContext.current

        Button(
            onClick = {
                Toast.makeText(context, contentDescription, Toast.LENGTH_SHORT)
                    .show()
            },
            colors = ButtonDefaults.buttonColors(
                containerColor = Color(0xFFE5DADA), // Timberwolf Gray
                contentColor = Color(0xFF02040F) // Rich Black
            ),
            modifier = Modifier
                .width(64.dp),
            contentPadding = PaddingValues(0.dp)
        ) {
            Icon(
                imageVector = imageVector,
                contentDescription = contentDescription,
                modifier = Modifier,
                tint = Color(0xFF02040F) // Rich Black
            )
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
    fun ActivityItem(activity: String, incomplete: Boolean = false) {
        val textColor = if (incomplete) Color(0xFF6E6E6E) else Color(0xFF02040F)
        val backgroundColor = if (incomplete) Color(0xFFDADADA) else Color(0xFFE5DADA)
        val fontStyle = if (incomplete) FontStyle.Italic else FontStyle.Normal
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
                text = activity,
                color = textColor,
                fontStyle = fontStyle,
                fontWeight = FontWeight.Normal,
                modifier = Modifier.weight(1f)
            )

            if (incomplete) {
                Icon(
                    imageVector = Icons.Default.Edit,
                    contentDescription = "Add activity",
                    modifier = Modifier
                        .padding(start = 8.dp)
                        .size(16.dp),
                    tint = textColor
                )
            }
        }
    }
}
