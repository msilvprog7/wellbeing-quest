package app.wellbeingquest.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.ForeignKey
import androidx.room.Index

@Entity(
    tableName = "entryFeelings",
    primaryKeys = ["entryId", "feelingId"],
    foreignKeys = [
        ForeignKey(
            entity = EntryEntity::class,
            parentColumns = ["id"],
            childColumns = ["entryId"],
            onDelete = ForeignKey.CASCADE,
            onUpdate = ForeignKey.CASCADE
        ),
        ForeignKey(
            entity = FeelingEntity::class,
            parentColumns = ["id"],
            childColumns = ["feelingId"],
            onDelete = ForeignKey.CASCADE,
            onUpdate = ForeignKey.CASCADE
        )
    ],
    indices = [
        Index(value = ["entryId"]),
        Index(value = ["feelingId"])
    ])
data class EntryFeelingEntity(
    @ColumnInfo(name = "entryId")
    val entryId: Int,

    @ColumnInfo(name = "feelingId")
    val feelingId: Int
)