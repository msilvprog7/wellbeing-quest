package app.wellbeingquest.ui.theme

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Close
import androidx.compose.material3.Icon
import androidx.compose.material3.InputChip
import androidx.compose.material3.Text
import androidx.compose.material3.TextField
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@Composable
fun GroupText(text: String, modifier: Modifier) {
    Text(
        text = text,
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
fun MultiEntryTextField(
    value: String,
    onValueChange: (String) -> Unit,
    label: @Composable () -> Unit,
    placeholder: @Composable () -> Unit,
    entries: List<String>,
    onDone: () -> Unit,
    onRemove: (String) -> Unit,
    enabled: Boolean = true,
) {
    Column(modifier = Modifier.padding(16.dp)) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier
                .fillMaxWidth()
                .padding(0.dp)
        ) {
            TextField(
                value = value,
                onValueChange = onValueChange,
                label = label,
                placeholder = placeholder,
                singleLine = true,
                keyboardOptions = KeyboardOptions.Default.copy(imeAction = ImeAction.Done),
                keyboardActions = KeyboardActions(onDone = {
                    onDone()
                })
            )

            NavigationButton(
                imageVector = Icons.Default.Add,
                contentDescription = "add feeling",
                onClick = {
                    onDone()
                },
                enabled = enabled,
                modifier = Modifier.padding(start = 8.dp)
            )
        }

        FlowRow(
            // mainAxisSpacing = 8.dp,
            // crossAxisSpacing = 8.dp
        ) {
            entries.forEach { entry ->
                InputChip(
                    selected = false,
                    onClick = {},
                    label = { Text(entry) },
                    trailingIcon = {
                        Icon(
                            imageVector = Icons.Default.Close,
                            contentDescription = "Remove",
                            modifier = Modifier.clickable {
                                onRemove(entry)
                            }
                        )
                    }
                )
            }


        }
    }
}