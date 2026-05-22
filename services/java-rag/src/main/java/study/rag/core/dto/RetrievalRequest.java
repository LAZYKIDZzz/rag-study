package study.rag.core.dto;

import jakarta.validation.constraints.Max;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.Map;

public record RetrievalRequest(
    @NotBlank String query,
    @JsonProperty("top_k") @Min(1) @Max(20) Integer topK,
    Map<String, Object> filters) {
}
