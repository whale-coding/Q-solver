<template>
  <div class="question-node" :class="{ selected: data.selected }">
    <div class="node-index">{{ data.index }}</div>
    <div class="node-content">
      <div class="node-title">{{ data.title }}</div>
    </div>
    <Handle type="target" :position="Position.Top" />
    <Handle type="source" :position="Position.Bottom" />
  </div>
</template>

<script setup>
import { Handle, Position } from '@vue-flow/core'

defineProps({
  data: {
    type: Object,
    required: true
  }
})
</script>

<style scoped>
.question-node {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px 8px 8px;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.7), rgba(79, 70, 229, 0.6));
  border: 2px solid rgba(99, 102, 241, 0.6);
  border-radius: 20px;
  min-width: 80px;
  max-width: 140px;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.question-node:hover {
  border-color: rgba(129, 140, 248, 1);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
  transform: translateY(-1px);
}

.question-node.selected {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.8), rgba(5, 150, 105, 0.7));
  border-color: rgba(16, 185, 129, 1);
  box-shadow: 0 4px 16px rgba(16, 185, 129, 0.4);
}

.node-index {
  width: 22px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  flex-shrink: 0;
}

.question-node.selected .node-index {
  background: rgba(255, 255, 255, 0.25);
}

.node-content {
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.node-title {
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Handle 样式 */
:deep(.vue-flow__handle) {
  width: 8px;
  height: 8px;
  background: rgba(99, 102, 241, 0.8);
  border: 2px solid rgba(255, 255, 255, 0.5);
  opacity: 0;
  transition: opacity 0.2s;
}

.question-node:hover :deep(.vue-flow__handle) {
  opacity: 1;
}

.question-node.selected :deep(.vue-flow__handle) {
  background: rgba(16, 185, 129, 0.8);
}
</style>
