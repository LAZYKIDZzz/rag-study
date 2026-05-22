package study.rag.core;

import org.springframework.stereotype.Service;
import study.rag.core.dto.AddFactRequest;
import study.rag.core.dto.CapabilitiesResponse;
import study.rag.core.dto.ChatMessageRecord;
import study.rag.core.dto.ChatMessageRequest;
import study.rag.core.dto.ChatMessageResponse;
import study.rag.core.dto.ChatSessionResponse;
import study.rag.core.dto.DocumentCreateRequest;
import study.rag.core.dto.DocumentListResponse;
import study.rag.core.dto.DocumentResponse;
import study.rag.core.dto.IndexResponse;
import study.rag.core.dto.MemoryResponse;
import study.rag.core.dto.RetrievalRequest;
import study.rag.core.dto.RetrievalResponse;
import study.rag.core.dto.RetrievedChunk;
import study.rag.core.model.ChatMessage;
import study.rag.core.model.ChatSession;
import study.rag.core.model.Chunk;
import study.rag.core.model.Document;
import study.rag.core.model.SearchResult;
import study.rag.core.model.UserMemory;

import java.time.Instant;
import java.util.List;
import java.util.Map;

@Service
public class RagService {
  private final DocumentRepository documents;
  private final ChatRepository chats;
  private final MemoryStore memory;
  private final TextPreprocessor preprocessor;
  private final ParagraphChunker chunker;
  private final LocalHashEmbeddingProvider embeddings;
  private final InMemoryVectorStore vectorStore;
  private final QueryProcessor queryProcessor;
  private final ExtractiveAnswerGenerator generator;

  public RagService(
      DocumentRepository documents,
      ChatRepository chats,
      MemoryStore memory,
      TextPreprocessor preprocessor,
      ParagraphChunker chunker,
      LocalHashEmbeddingProvider embeddings,
      InMemoryVectorStore vectorStore,
      QueryProcessor queryProcessor,
      ExtractiveAnswerGenerator generator) {
    this.documents = documents;
    this.chats = chats;
    this.memory = memory;
    this.preprocessor = preprocessor;
    this.chunker = chunker;
    this.embeddings = embeddings;
    this.vectorStore = vectorStore;
    this.queryProcessor = queryProcessor;
    this.generator = generator;
  }

  public CapabilitiesResponse capabilities() {
    return new CapabilitiesResponse(
        "java-rag",
        List.of(
            "document-ingestion",
            "preprocessing",
            "chunking",
            "local-embeddings",
            "in-memory-vector-search",
            "query-rewriting",
            "chat-memory",
            "extractive-answer-generation"),
        "local-hash-embedding-128d",
        "in-memory-cosine",
        "in-memory");
  }

  public DocumentResponse createDocument(DocumentCreateRequest request) {
    Document document = documents.create(
        request.title().trim(),
        request.content(),
        request.metadata() == null ? Map.of() : request.metadata());
    return documentResponse(document);
  }

  public DocumentListResponse listDocuments() {
    return new DocumentListResponse(documents.list().stream().map(this::documentResponse).toList());
  }

  public IndexResponse indexDocument(String documentId) {
    Document document = documents.get(documentId);
    String cleanText = preprocessor.clean(document.content());
    List<Chunk> chunks = chunker.chunk(document, cleanText);
    List<double[]> vectors = chunks.stream().map(chunk -> embeddings.embed(chunk.content())).toList();
    vectorStore.deleteByDocument(documentId);
    vectorStore.upsertMany(chunks, vectors);
    documents.replaceChunks(documentId, chunks);
    return new IndexResponse(documentId, chunks.size(), vectorStore.countForDocument(documentId));
  }

  public RetrievalResponse search(RetrievalRequest request) {
    String rewrittenQuery = queryProcessor.rewrite(request.query(), null, memory.get("anonymous"));
    List<SearchResult> results = retrieve(
        rewrittenQuery,
        request.topK() == null ? 5 : request.topK(),
        request.filters() == null ? Map.of() : request.filters());
    return new RetrievalResponse(
        request.query(),
        rewrittenQuery,
        results.stream().map(this::retrievedChunk).toList());
  }

  public ChatSessionResponse createSession(String userId) {
    return sessionResponse(chats.create(userId));
  }

  public ChatSessionResponse getSession(String sessionId) {
    return sessionResponse(chats.get(sessionId));
  }

  public ChatMessageResponse sendMessage(String sessionId, ChatMessageRequest request) {
    ChatSession session = chats.get(sessionId);
    String userId = request.userId() == null || request.userId().isBlank() ? session.userId() : request.userId();
    UserMemory userMemory = memory.get(userId);
    String rewrittenQuery = queryProcessor.rewrite(request.message(), session, userMemory);
    List<SearchResult> results = retrieve(rewrittenQuery, request.topK() == null ? 5 : request.topK(), Map.of());
    String answer = generator.generate(request.message(), results, userMemory);

    chats.addMessage(sessionId, new ChatMessage("user", request.message(), Instant.now(), List.of()));
    chats.addMessage(sessionId, new ChatMessage("assistant", answer, Instant.now(), results));
    MemoryResponse updatedMemory = memoryResponse(memory.updateSummary(userId, request.message()));

    return new ChatMessageResponse(
        sessionId,
        answer,
        rewrittenQuery,
        results.stream().map(this::retrievedChunk).toList(),
        results.stream().map(this::retrievedChunk).toList(),
        List.of(updatedMemory.recentSummary()));
  }

  public MemoryResponse getMemory(String userId) {
    return memoryResponse(memory.get(userId));
  }

  public MemoryResponse addFact(String userId, String fact) {
    return memoryResponse(memory.addFact(userId, fact));
  }

  private List<SearchResult> retrieve(String query, int topK, Map<String, Object> filters) {
    return vectorStore.search(embeddings.embed(query), topK, filters);
  }

  private DocumentResponse documentResponse(Document document) {
    return new DocumentResponse(
        document.id(),
        document.title(),
        document.metadata(),
        document.indexed(),
        documents.chunkCount(document.id()),
        document.createdAt());
  }

  private RetrievedChunk retrievedChunk(SearchResult result) {
    return new RetrievedChunk(
        result.chunk().id(),
        result.chunk().documentId(),
        result.chunk().title(),
        result.chunk().ordinal(),
        result.chunk().content(),
        Math.round(result.score() * 1_000_000.0) / 1_000_000.0,
        result.chunk().metadata());
  }

  private ChatSessionResponse sessionResponse(ChatSession session) {
    return new ChatSessionResponse(
        session.id(),
        session.userId(),
        session.messages().stream()
            .map(message -> new ChatMessageRecord(
                message.role(),
                message.content(),
                message.createdAt(),
                message.citations().stream().map(this::retrievedChunk).toList()))
            .toList(),
        session.createdAt());
  }

  private MemoryResponse memoryResponse(UserMemory userMemory) {
    return new MemoryResponse(userMemory.userId(), userMemory.facts(), userMemory.recentSummary());
  }
}
