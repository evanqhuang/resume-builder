import type { DragEndEvent } from '@dnd-kit/core';
import { arrayMove } from '@dnd-kit/sortable';
import { useResume } from './useResume';
import { saveOrder } from '../services/api';
import type { SectionName } from '../types/resume';

export function useSectionReorder(section: SectionName, items: { id: string }[]) {
  const { dispatch } = useResume();

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIdx = items.findIndex(item => item.id === active.id);
    const newIdx = items.findIndex(item => item.id === over.id);
    const newOrder = arrayMove(items.map(item => item.id), oldIdx, newIdx);
    const previousOrder = items.map(item => item.id);

    dispatch({ type: 'REORDER_SECTION', payload: { section, order: newOrder } });

    try {
      await saveOrder({ [section]: newOrder });
    } catch (error) {
      console.error(`Failed to save ${section} order:`, error);
      dispatch({ type: 'REORDER_SECTION', payload: { section, order: previousOrder } });
    }
  };

  return { handleDragEnd };
}
