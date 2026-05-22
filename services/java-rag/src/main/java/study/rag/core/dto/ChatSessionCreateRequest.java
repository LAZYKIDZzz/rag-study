package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

public record ChatSessionCreateRequest(@JsonProperty("user_id") String userId, String title) {
}
