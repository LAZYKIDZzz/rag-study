package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record CapabilitiesResponse(
    String service,
    List<String> features,
    @JsonProperty("embedding_provider") String embeddingProvider,
    @JsonProperty("vector_store") String vectorStore,
    @JsonProperty("memory_strategy") String memoryStrategy) {
}
