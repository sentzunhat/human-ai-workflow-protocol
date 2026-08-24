# Troubleshooting Context Reshaping

## Quick Diagnostics

```bash
# Check HAWP version
hawp --version  # Should show v0.0.3+

# Check if Ollama is running
curl -s http://localhost:11434/api/version && echo "✓ Ollama running" || echo "✗ Ollama not found"

# List available Ollama models
ollama list

# Test Ollama embedding
curl -X POST http://localhost:11434/api/embeddings -H "Content-Type: application/json" \
  -d '{"model":"nomic-embed-text","prompt":"test"}' | head -20

# Check HAWP config location
ls -la ~/.hawp/config/context.json
```

---

## Common Issues

### 1. "Ollama not found" or "Connection refused"

**Error Message:**
```
error: failed to initialize LLM: dial tcp [::1]:11434: connect: connection refused
```

**Cause:** `--llm-reshape` requires Ollama running locally.

**Solution:**
```bash
# Step 1: Install Ollama
# Visit: https://ollama.ai/download

# Step 2: Start Ollama server (in background or separate terminal)
ollama serve

# Step 3: Download a model (required before first use)
ollama pull mistral

# Step 4: Try the command again
hawp search "query" --context --llm-reshape
```

**Verification:**
```bash
# Check if Ollama is running
curl -s http://localhost:11434/api/version
# Should output: {"version":"0.1.45"} (or similar)

# Check if mistral model is downloaded
ollama list | grep mistral
# Should show: mistral     latest     3b87...   4.0 GB
```

---

### 2. "Model not found"

**Error Message:**
```
error: failed to reshape: model 'mistral' not found
```

**Cause:** Ollama model hasn't been downloaded yet.

**Solution:**
```bash
# List downloaded models
ollama list

# Download the missing model
ollama pull mistral

# Verify it's there
ollama list | grep mistral

# Try command again
hawp search "query" --context --llm-reshape
```

**Popular Models:**
```bash
# Fastest (recommended for first try)
ollama pull neural-chat

# Best balance (recommended)
ollama pull mistral

# Highest quality (slowest, requires GPU ideally)
ollama pull llama2-13b

# For embeddings only
ollama pull nomic-embed-text
ollama pull mxbai-embed-large
```

---

### 3. ONNX Model Download Slow

**Problem:** First run takes 30-60 seconds downloading embeddings model.

**Error/Output:**
```
Downloading ONNX model: all-MiniLM-L6-v2... (this may take a minute)
```

**Why:** ONNX models (~50MB) are downloaded on first use and cached for future runs.

**Solution:**
```bash
# Pre-download models before running searches
go run ./cmd/hawp search "dummy" --context 2>&1 | head -10

# Or wait - it only happens once
# Future runs use cached model (instant)

# Verify cache location
ls -la ~/.cache/hawp/onnx/  # On Linux/Mac
ls -la %APPDATA%/hawp/onnx/  # On Windows
```

**Troubleshooting:**
```bash
# Check disk space (needs ~100MB free)
df -h | grep $HOME

# Clear cache if needed
rm -rf ~/.cache/hawp/onnx/

# Try again
hawp search "query" --context --llm-reshape
```

---

### 4. Slow Performance

**Problem:** Reshaping takes >10 seconds.

**Factors:**
- CPU-only (no GPU acceleration)
- Large LLM model (mistral 7B vs neural-chat)
- Long context (many search results)
- Disk I/O (network latency)

**Solutions (in order of impact):**

1. **Use faster LLM model:**
   ```json
   {
     "llm": {
       "model": "neural-chat:latest"  // 30% faster than mistral
     }
   }
   ```
   Test: `hawp search "test" --context --llm-reshape --verbose`

2. **Reduce token budget:**
   ```json
   {
     "llm": {
       "max_tokens": 256  // Default is 512
     }
   }
   ```
   Tradeoff: Less detailed reshaping.

3. **Reduce key concepts:**
   ```json
   {
     "reshaping": {
       "top_k": 3  // Default is 5
     }
   }
   ```
   Tradeoff: Fewer themes identified.

4. **Use faster embeddings (ONNX > Ollama):**
   ```json
   {
     "embeddings": {
       "backend": "onnx",
       "model": "all-MiniLM-L6-v2"  // ~40ms
     }
   }
   ```
   vs
   ```json
   {
     "embeddings": {
       "backend": "ollama",
       "model": "nomic-embed-text"  // ~150ms (but higher quality)
     }
   }
   ```

5. **Enable GPU acceleration:**
   - Install CUDA (NVIDIA) or Metal (Mac)
   - Ollama auto-detects and uses GPU if available
   - 5-10x speedup if GPU available

**Performance Profile (Mac CPU-only):**
```
Neural-chat + ONNX embeddings: ~1.2s ✓ Recommended
Mistral + ONNX embeddings:     ~1.8s ✓ Good
Llama2-13b + ONNX:            ~3-5s ⚠️ Slow
Llama2-70b + ONNX:            ~10+ s ✗ Very slow
```

---

### 5. Memory Issues

**Error Message:**
```
runtime: out of memory
```

**Cause:** Large models use significant RAM (7-40GB depending on model).

**Solution:**
```bash
# Check available RAM
free -h  # Linux
vm_stat  # Mac (look for "Pages free")

# Use smaller model
ollama pull neural-chat:latest  # ~3GB
# vs
ollama pull llama2-70b          # ~40GB
```

**If still out of memory:**
1. Close other applications
2. Reduce token budget (max_tokens=256)
3. Use CPU-only ONNX embeddings (smaller overhead)

---

### 6. Network Timeout Issues

**Error Message:**
```
context deadline exceeded
```

**Cause:** Ollama HTTP request timed out (default 30s for embeddings, 60s for LLM).

**Solution:**
```bash
# Check Ollama status
curl -s http://localhost:11434/api/version

# Restart Ollama
killall ollama
ollama serve

# Increase timeout in config
{
  "reshaping": {
    "timeout_ms": 60000  // 60 seconds (default 30)
  }
}

# Try again
hawp search "query" --context --llm-reshape
```

---

### 7. ONNX Model Download Failure

**Error Message:**
```
error: failed to download ONNX model: 404 Not Found
```

**Cause:** Network issue or model unavailable.

**Solution:**
```bash
# Check internet connection
ping google.com

# Try manual download
go run ./cmd/hawp search "dummy" --context -vv

# Check logs for download URL
# Try again with better connection

# Use cached version if available
ls -la ~/.cache/hawp/onnx/
```

---

### 8. Configuration Not Applied

**Problem:** Changes to `~/.hawp/config/context.json` not taking effect.

**Solution:**
```bash
# Verify config format
cat ~/.hawp/config/context.json | jq .  # Should be valid JSON

# Restart HAWP (no daemon, always fresh)
hawp search "query" --context --llm-reshape

# Or use environment variables (highest priority)
export HAWP_EMBEDDINGS_BACKEND=ollama
export HAWP_LLM_BACKEND=ollama
hawp search "query" --context --llm-reshape
```

**Config Priority:**
1. Environment variables (highest)
2. `~/.hawp/config/context.json`
3. Hardcoded defaults (lowest)

---

### 9. Ollama Model Doesn't Exist

**Error Message:**
```
error: failed to reshape: model 'quantum-brain-9000' not found
```

**Cause:** Model name typo or model not downloaded.

**Solution:**
```bash
# List available models
ollama list

# Popular models
# For LLM: mistral, neural-chat, llama2, dolphin-mixtral
# For embeddings: nomic-embed-text, mxbai-embed-large, all-minilm

# Download if needed
ollama pull mistral

# Verify and try again
hawp search "query" --context --llm-reshape
```

---

### 10. HAWP Works, But --llm-reshape Crashes

**Error Message:**
```
error: unexpected error: <some panic>
```

**Debugging:**
```bash
# Run with verbose logging
hawp search "query" --context --llm-reshape -vv

# Get full stack trace
hawp search "query" --context --llm-reshape -vv 2>&1 > debug.log
cat debug.log | head -100

# Check if Ollama is still running
curl -s http://localhost:11434/api/version

# Restart both
killall ollama
ollama serve

# Try again
hawp search "query" --context --llm-reshape
```

**If still failing:**
1. Check GitHub issues: https://github.com/sentzunhat/hawp/issues
2. Provide debug.log output
3. Include: OS, HAWP version, Ollama version, model names

---

### 11. High CPU Usage During Reshaping

**Observation:** CPU at 100%, takes a long time.

**Expected:** Yes, this is normal during LLM inference.

**What's happening:**
- ONNX embeddings run locally on CPU
- Ollama LLM also runs locally on CPU
- No GPU acceleration available

**To optimize:**
1. **Install GPU drivers** (dramatic speedup if available):
   - NVIDIA: Install CUDA
   - Apple Silicon: Metal (auto-enabled)
   - AMD: Install ROCm

2. **Or accept slowdown:**
   - CPU-only is expected behavior
   - Embeddings: 30-50ms
   - LLM: 1-3 seconds per reshape
   - Total: 1.5-3.5 seconds per search

---

### 12. Ollama Port Already in Use

**Error Message:**
```
listen tcp [::1]:11434: bind: address already in use
```

**Cause:** Ollama already running (or different service on port 11434).

**Solution:**
```bash
# Find process using port 11434
lsof -i :11434  # Mac/Linux
netstat -ano | findstr :11434  # Windows

# Kill existing Ollama
killall ollama
sleep 2

# Restart
ollama serve

# Or use different port (requires config update)
# Not supported yet in v0.0.3, use localhost:11434
```

---

## Debug Checklist

Use this checklist when something doesn't work:

```
[ ] HAWP version >= v0.0.3: hawp --version
[ ] Ollama running: curl http://localhost:11434/api/version
[ ] Ollama model exists: ollama list | grep <model>
[ ] Disk space available: df -h | grep $HOME
[ ] Internet connected: ping google.com
[ ] Config valid JSON: cat ~/.hawp/config/context.json | jq .
[ ] Try verbose logging: hawp search "query" --context --llm-reshape -vv
[ ] Check GitHub issues: https://github.com/sentzunhat/hawp/issues
```

---

## Still Stuck?

### Minimal Repro
```bash
# 1. Start fresh
killall ollama
ollama serve  # New terminal

# 2. Download model
ollama pull mistral

# 3. Test Ollama directly
curl -X POST http://localhost:11434/api/embeddings \
  -H "Content-Type: application/json" \
  -d '{"model":"nomic-embed-text","prompt":"test"}'

# 4. Try HAWP
hawp search "test" --context --llm-reshape -vv

# 5. Capture output
hawp search "test" --context --llm-reshape -vv 2>&1 | head -50 > bug.txt
```

### Report Issue
Provide:
1. Output of `hawp --version`
2. Output of `ollama list`
3. Full error from `bug.txt` above
4. OS and architecture (Mac Apple Silicon, Linux x86_64, etc)

Link: https://github.com/sentzunhat/hawp/issues/new

---

## Performance Reference

### Expected Times (Mac CPU-only, 5 results, ~2000 tokens)

| Embeddings | LLM | Embed Time | LLM Time | Total |
|---|---|---|---|---|
| ONNX (384-dim) | neural-chat | 40ms | 600ms | ~650ms |
| ONNX (768-dim) | neural-chat | 45ms | 600ms | ~650ms |
| ONNX (384-dim) | mistral | 40ms | 1200ms | ~1.2s |
| Ollama | neural-chat | 150ms | 600ms | ~750ms |
| Ollama | mistral | 150ms | 1200ms | ~1.4s |

**With GPU (NVIDIA):** 5-10x faster  
**With GPU (Mac Silicon M1/M2):** 3-5x faster

---

## Common Success Patterns

**Pattern 1: Quick Start**
```bash
ollama pull mistral
ollama serve
hawp search "kubernetes" --context --llm-reshape
```
Expected: <2 seconds, good quality results

**Pattern 2: All-Ollama**
```bash
ollama pull nomic-embed-text
ollama pull mistral
ollama serve
# Config: embeddings.backend="ollama", llm.backend="ollama"
hawp search "kubernetes" --context --llm-reshape
```
Expected: ~1.5 seconds, highest quality, single service

**Pattern 3: Offline-First (ONNX)**
```bash
# First run (downloads model)
hawp search "kubernetes" --context --llm-reshape
# ~30s (includes model download)

# Future runs (cached)
hawp search "kubernetes" --context --llm-reshape
# <2s
```
Expected: ONNX works offline, Ollama requires network/service
