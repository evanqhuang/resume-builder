interface KeywordBadgesProps {
  keywords: string[];
}

export const KeywordBadges = ({ keywords }: KeywordBadgesProps) => {
  if (keywords.length === 0) return null;

  return (
    <div className="space-y-2">
      <h3 className="text-sm font-semibold text-gray-700">Keywords Detected</h3>
      <div className="flex flex-wrap gap-2">
        {keywords.map((keyword, index) => (
          <span
            key={index}
            className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800"
          >
            {keyword}
          </span>
        ))}
      </div>
    </div>
  );
};
