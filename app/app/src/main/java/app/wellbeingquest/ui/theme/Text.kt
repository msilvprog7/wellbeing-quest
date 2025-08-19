package app.wellbeingquest.ui.theme

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.FlowRow
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Add
import androidx.compose.material.icons.filled.Close
import androidx.compose.material3.DropdownMenu
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.InputChip
import androidx.compose.material3.MenuAnchorType
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
import androidx.compose.ui.window.PopupProperties

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

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AutoCompleteTextField(
    text: String,
    onTextChange: (String) -> Unit,
    label: @Composable () -> Unit,
    placeholder: @Composable () -> Unit,
    suggestions: List<String>,
    onSuggestionsChange: (List<String>) -> Unit,
    expanded: Boolean,
    onExpandedChange: (Boolean) -> Unit,
    modifier: Modifier = Modifier,
    onDone: (String) -> Unit = {}
) {

    ExposedDropdownMenuBox(
        expanded = expanded,
        onExpandedChange = onExpandedChange
    ) {
        TextField(
            value = text,
            onValueChange = {
                onTextChange(it)
            },
            label = label,
            placeholder = placeholder,
            modifier = modifier
                .menuAnchor(MenuAnchorType.PrimaryNotEditable, true),
            singleLine = true,
            trailingIcon = {
                ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded)
            },
            keyboardOptions = KeyboardOptions.Default.copy(imeAction = ImeAction.Done),
            keyboardActions = KeyboardActions(onDone = {
                onDone(text)
            })
        )

        DropdownMenu(
            expanded = expanded,
            onDismissRequest = { onExpandedChange(false) },
            properties = PopupProperties(focusable = false)
        ) {
            suggestions.forEach { suggestion ->
                DropdownMenuItem(
                    text = { Text(suggestion) },
                    onClick = {
                        onTextChange(suggestion)
                        onExpandedChange(false)
                        onDone(suggestion)
                    }
                )
            }
        }
    }
}

@Composable
fun MultiEntryTextField(
    text: String,
    onTextChange: (String) -> Unit,
    label: @Composable () -> Unit,
    placeholder: @Composable () -> Unit,
    suggestions: List<String>,
    onSuggestionsChange: (List<String>) -> Unit,
    expanded: Boolean,
    onExpandedChange: (Boolean) -> Unit,
    selected: List<String>,
    onSelectedChange: (List<String>) -> Unit,
    modifier: Modifier = Modifier
) {
    var onDone: (String) -> Unit = {
        if (it.isNotBlank() && !selected.contains(it.trim())) {
            // Add to list
            onSelectedChange(selected + it.trim())

            // Empty text
            onTextChange("")
        }
    }

    Column(modifier = Modifier.padding(16.dp)) {
        Row(
            verticalAlignment = Alignment.CenterVertically,
            modifier = Modifier
                .fillMaxWidth()
                .padding(0.dp)
        ) {
            AutoCompleteTextField(
                text = text,
                onTextChange = onTextChange,
                label = label,
                placeholder = placeholder,
                suggestions = suggestions.filterNot { selected.contains(it) },
                onSuggestionsChange = onSuggestionsChange,
                expanded = expanded,
                onExpandedChange = onExpandedChange,
                modifier = modifier,
                onDone = onDone
            )

            NavigationButton(
                imageVector = Icons.Default.Add,
                contentDescription = "add feeling",
                onClick = {
                    onDone(text)
                },
                enabled = text.isNotBlank() && !selected.contains(text.trim()),
                modifier = Modifier.padding(start = 8.dp)
            )
        }

        FlowRow {
            selected.forEach { entry ->
                InputChip(
                    selected = false,
                    onClick = {},
                    label = { Text(entry) },
                    trailingIcon = {
                        Icon(
                            imageVector = Icons.Default.Close,
                            contentDescription = "Remove",
                            modifier = Modifier.clickable {
                                // Remove
                                onSelectedChange(selected - entry)
                            }
                        )
                    }
                )
            }


        }
    }
}