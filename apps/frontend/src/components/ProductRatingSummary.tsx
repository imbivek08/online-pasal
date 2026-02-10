import type { ProductRatingStats } from '../lib/api';
import RatingStars from './RatingStars';

interface ProductRatingSummaryProps {
  stats: ProductRatingStats;
}

export default function ProductRatingSummary({ stats }: ProductRatingSummaryProps) {
  const ratingBreakdown = [
    { stars: 5, count: stats.five_star_count },
    { stars: 4, count: stats.four_star_count },
    { stars: 3, count: stats.three_star_count },
    { stars: 2, count: stats.two_star_count },
    { stars: 1, count: stats.one_star_count },
  ];

  const getPercentage = (count: number) => {
    if (stats.total_reviews === 0) return 0;
    return (count / stats.total_reviews) * 100;
  };

  return (
    <div className="bg-gray-50 rounded-lg p-6 mb-8">
      <h3 className="text-lg font-semibold mb-4">Customer Reviews</h3>
      
      <div className="flex items-start gap-8">
        {/* Overall Rating */}
        <div className="text-center">
          <div className="text-4xl font-bold text-gray-900 mb-2">
            {stats.average_rating.toFixed(1)}
          </div>
          <RatingStars rating={stats.average_rating} size="md" />
          <p className="text-sm text-gray-600 mt-2">
            {stats.total_reviews} {stats.total_reviews === 1 ? 'review' : 'reviews'}
          </p>
        </div>

        {/* Rating Breakdown */}
        <div className="flex-1">
          {ratingBreakdown.map(({ stars, count }) => (
            <div key={stars} className="flex items-center gap-3 mb-2">
              <span className="text-sm text-gray-600 w-12">{stars} star</span>
              <div className="flex-1 bg-gray-200 rounded-full h-2">
                <div
                  className="bg-yellow-400 h-2 rounded-full transition-all"
                  style={{ width: `${getPercentage(count)}%` }}
                />
              </div>
              <span className="text-sm text-gray-600 w-12 text-right">{count}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
