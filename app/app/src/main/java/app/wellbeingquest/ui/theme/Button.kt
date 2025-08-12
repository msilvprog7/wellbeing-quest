package app.wellbeingquest.ui.theme

import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.width
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Icon
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.unit.dp

@Composable
fun NavigationButton(imageVector: ImageVector, contentDescription: String, onClick: () -> Unit, enabled: Boolean = true) {
    var context = LocalContext.current

    Button(
        onClick = onClick,
        enabled = enabled,
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
