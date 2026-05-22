package study.rag.core;

import org.junit.jupiter.api.Test;
import study.rag.core.dto.ChatMessageRequest;
import study.rag.core.dto.DocumentCreateRequest;
import study.rag.core.dto.RetrievalRequest;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;

class RagServiceTest {
  @Test
  void runsFullRagFlow() {
    RagService service = new RagService(
        new DocumentRepository(),
        new ChatRepository(),
        new MemoryStore(),
        new TextPreprocessor(),
        new ParagraphChunker(),
        new LocalHashEmbeddingProvider(),
        new InMemoryVectorStore(),
        new QueryProcessor(),
        new ExtractiveAnswerGenerator());

    var document = service.createDocument(new DocumentCreateRequest(
        "RAG basics",
        "Chunking splits documents. Embeddings convert text into vectors for retrieval.",
        Map.of("source", "test")));
    service.indexDocument(document.id());

    var search = service.search(new RetrievalRequest("What are embeddings?", 3, Map.of()));
    assertThat(search.retrievedChunks()).isNotEmpty();

    var session = service.createSession("demo-user");
    service.addFact("demo-user", "The user is studying RAG.");
    var response = service.sendMessage(session.id(), new ChatMessageRequest(
        "demo-user",
        "What does chunking do?",
        3));

    assertThat(response.answer()).contains("retrieved knowledge");
    assertThat(response.citations()).isNotEmpty();
    assertThat(response.memoryUpdates()).isNotEmpty();
  }
}
