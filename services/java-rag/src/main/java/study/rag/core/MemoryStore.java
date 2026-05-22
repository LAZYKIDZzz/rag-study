package study.rag.core;

import org.springframework.stereotype.Repository;
import study.rag.core.model.UserMemory;

import java.util.HashMap;
import java.util.Map;

@Repository
public class MemoryStore {
  private final Map<String, UserMemory> memories = new HashMap<>();

  public synchronized UserMemory get(String userId) {
    return memories.computeIfAbsent(userId, UserMemory::new);
  }

  public synchronized UserMemory addFact(String userId, String fact) {
    UserMemory memory = get(userId);
    memory.addFact(fact.trim());
    return memory;
  }

  public synchronized UserMemory updateSummary(String userId, String latestQuestion) {
    UserMemory memory = get(userId);
    memory.updateSummary(latestQuestion);
    return memory;
  }
}
