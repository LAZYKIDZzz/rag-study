package study.rag.core.dto;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public record ChatMessageResponse(
    @JsonProperty("session_id") String sessionId,
    String answer,
    @JsonProperty("rewritten_query") String rewrittenQuery,
    List<RetrievedChunk> citations,
    @JsonProperty("retrieved_chunks") List<RetrievedChunk> retrievedChunks,
    @JsonProperty("memory_updates") List<String> memoryUpdates) {
}
