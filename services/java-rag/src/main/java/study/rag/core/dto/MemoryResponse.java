package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record MemoryResponse(
    @JsonProperty("user_id") String userId,
    List<String> facts,
    @JsonProperty("recent_summary") String recentSummary) {
}
