import { formatDistanceToNow } from 'date-fns';
import { ThumbsUp, CheckCircle } from 'lucide-react';
import RatingStars from './RatingStars';
import type { Review } from '../lib/api';

interface ReviewCardProps {
  review: Review;
  onMarkHelpful?: (reviewId: string) => void;
}

export default function ReviewCard({ review, onMarkHelpful }: ReviewCardProps) {
  return (
    <div className="border-b border-gray-200 pb-6 mb-6 last:border-b-0">
      <div className="flex items-start gap-4">
        {/* User Avatar */}
        <div className="flex-shrink-0">
          {review.user_avatar ? (
            <img
              src={review.user_avatar}
              alt={review.user_name || 'User'}
              className="w-10 h-10 rounded-full"
            />
          ) : (
            <div className="w-10 h-10 rounded-full bg-gradient-to-br from-purple-400 to-pink-400 flex items-center justify-center text-white font-semibold">
              {review.user_name?.charAt(0) || 'U'}
            </div>
          )}
        </div>

        {/* Review Content */}
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <span className="font-semibold text-gray-900">
              {review.user_name || 'Anonymous'}
            </span>
            {review.is_verified_purchase && (
              <span className="inline-flex items-center gap-1 text-xs text-green-600 bg-green-50 px-2 py-1 rounded-full">
                <CheckCircle className="w-3 h-3" />
                Verified Purchase
              </span>
            )}
          </div>

          {/* Rating */}
          <div className="flex items-center gap-3 mb-2">
            <RatingStars rating={review.rating} size="sm" />
            <span className="text-sm text-gray-500">
              {formatDistanceToNow(new Date(review.created_at), { addSuffix: true })}
            </span>
          </div>

          {/* Title */}
          {review.title && (
            <h4 className="font-semibold text-gray-900 mb-2">{review.title}</h4>
          )}

          {/* Comment */}
          {review.comment && (
            <p className="text-gray-700 mb-3 leading-relaxed">{review.comment}</p>
          )}

          {/* Helpful Button */}
          <button
            onClick={() => onMarkHelpful?.(review.id)}
            className="inline-flex items-center gap-2 text-sm text-gray-600 hover:text-gray-900 transition-colors"
          >
            <ThumbsUp className="w-4 h-4" />
            Helpful {review.helpful_count > 0 && `(${review.helpful_count})`}
          </button>
        </div>
      </div>
    </div>
  );
}
