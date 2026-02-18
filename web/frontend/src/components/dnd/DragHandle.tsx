import { GripVertical } from 'lucide-react';
import type { DragHandleProps } from './SortableItem';

export const DragHandle = ({ listeners, attributes }: DragHandleProps) => (
  <button
    type="button"
    {...listeners}
    {...attributes}
    className="cursor-grab active:cursor-grabbing p-1 text-gray-400 hover:text-gray-600 touch-none"
  >
    <GripVertical className="w-4 h-4" />
  </button>
);
