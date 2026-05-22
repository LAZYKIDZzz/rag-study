package study.rag.core.dto;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

public record ChatMessageRequest(
    @JsonProperty("user_id") String userId,
    @NotBlank String message,
    @JsonProperty("top_k") @Min(1) @Max(20) Integer topK) {
}
