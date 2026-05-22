package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.Instant;
import java.util.Map;

public record DocumentResponse(
    String id,
    String title,
    Map<String, Object> metadata,
    boolean indexed,
    @JsonProperty("chunk_count") int chunkCount,
    @JsonProperty("created_at") Instant createdAt) {
}
