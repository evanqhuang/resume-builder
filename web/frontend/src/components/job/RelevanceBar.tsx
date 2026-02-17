interface RelevanceBarProps {
  score: number;
  showLabel?: boolean;
}

export const RelevanceBar = ({ score, showLabel = true }: RelevanceBarProps) => {
  const getColor = (score: number) => {
    if (score >= 70) return 'bg-green-500';
    if (score >= 40) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  const getTextColor = (score: number) => {
    if (score >= 70) return 'text-green-700';
    if (score >= 40) return 'text-yellow-700';
    return 'text-red-700';
  };

  return (
    <div className="flex items-center gap-2 min-w-[80px]">
      <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
        <div
          className={`h-full ${getColor(score)} transition-all duration-300`}
          style={{ width: `${score}%` }}
        />
      </div>
      {showLabel && (
        <span className={`text-xs font-medium ${getTextColor(score)} min-w-[32px]`}>
          {Math.round(score)}
        </span>
      )}
    </div>
  );
};
