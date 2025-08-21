package app.wellbeingquest.data.local.entity

import androidx.room.TypeConverter

class Converters {
    @TypeConverter
    fun fromStringList(value: List<String>): String {
        return value.joinToString(separator = ";")
    }

    @TypeConverter
    fun toStringList(value: String): List<String> {
        return if (value.isEmpty()) emptyList() else value.split(";")
    }
}