package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.List;

public record ChatMessageRecord(
    String role,
    String content,
    @JsonProperty("created_at") Instant createdAt,
    List<RetrievedChunk> citations) {
}
