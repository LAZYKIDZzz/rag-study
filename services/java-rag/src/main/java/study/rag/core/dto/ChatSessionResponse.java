package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

public record ChatSessionResponse(
    String id,
    @JsonProperty("user_id") String userId,
    List<ChatMessageRecord> messages,
    @JsonProperty("created_at") Instant createdAt) {
}
