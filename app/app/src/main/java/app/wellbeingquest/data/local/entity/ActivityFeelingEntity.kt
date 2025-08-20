package app.wellbeingquest.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.ForeignKey
import androidx.room.Index

@Entity(
    tableName = "activityFeelings",
    primaryKeys = ["activityId", "feelingId"],
    foreignKeys = [
        ForeignKey(
            entity = ActivityEntity::class,
            parentColumns = ["id"],
            childColumns = ["activityId"],
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
        Index(value = ["activityId"]),
        Index(value = ["feelingId"])
    ])
data class ActivityFeelingEntity(
    @ColumnInfo(name = "activityId")
    val activityId: Int,

    @ColumnInfo(name = "feelingId")
    val feelingId: Int
)