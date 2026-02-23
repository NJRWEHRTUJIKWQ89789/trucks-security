"use client";

import { useEffect, useState } from "react";
import { Star, Loader2 } from "lucide-react";
import { gql } from "@/lib/graphql";

interface Feedback {
  id: string;
  clientName: string;
  rating: number;
  comment: string;
  category: string;
  createdAt: string;
}

function StarRow({ rating, size = "w-4 h-4" }: { rating: number; size?: string }) {
  return (
    <div className="flex items-center gap-0.5">
      {Array.from({ length: 5 }).map((_, i) => (
        <Star
          key={i}
          className={`${size} ${
            i < rating ? "fill-yellow-400 text-yellow-400" : "text-gray-300"
          }`}
        />
      ))}
    </div>
  );
}

function formatDate(raw: string): string {
  try {
    const d = new Date(raw);
    return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
  } catch {
    return raw;
  }
}

export default function ClientFeedbackPage() {
  const [feedbacks, setFeedbacks] = useState<Feedback[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    gql<{ feedbacks: { items: Feedback[]; totalCount: number } }>(
      `{ feedbacks(page:1,perPage:50) { items { id clientName rating comment category createdAt } totalCount } }`
    )
      .then((d) => setFeedbacks(d.feedbacks.items))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  // Compute rating summary from real data
  const totalCount = feedbacks.length;
  const avgRating = totalCount > 0
    ? feedbacks.reduce((sum, f) => sum + f.rating, 0) / totalCount
    : 0;
  const roundedAvg = Math.round(avgRating * 10) / 10;

  const ratingBreakdown = [5, 4, 3, 2, 1].map((stars) => {
    const count = feedbacks.filter((f) => f.rating === stars).length;
    const percent = totalCount > 0 ? (count / totalCount) * 100 : 0;
    return { stars, count, percent };
  });

  if (loading) {
    return (
      <div className="flex justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">Client Feedback</h1>

      <div className="bg-white rounded-xl border border-gray-200 p-8 mb-8">
        <div className="flex flex-col md:flex-row items-start md:items-center gap-8">
          <div className="text-center">
            <div className="text-5xl font-bold text-gray-900">{roundedAvg.toFixed(1)}</div>
            <div className="text-sm text-gray-500 mt-1">out of 5.0</div>
            <StarRow rating={Math.round(avgRating)} size="w-5 h-5" />
            <div className="text-sm text-gray-500 mt-2">{totalCount} reviews</div>
          </div>

          <div className="flex-1 w-full space-y-2">
            {ratingBreakdown.map((row) => (
              <div key={row.stars} className="flex items-center gap-3">
                <span className="text-sm font-medium text-gray-600 w-12">{row.stars} star</span>
                <div className="flex-1 h-3 bg-gray-100 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-yellow-400 rounded-full"
                    style={{ width: `${row.percent}%` }}
                  />
                </div>
                <span className="text-sm text-gray-500 w-10 text-right">{row.count}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="space-y-4">
        {feedbacks.map((review) => (
          <div
            key={review.id}
            className="bg-white rounded-xl border border-gray-200 p-6"
          >
            <div className="flex items-start justify-between mb-3">
              <div>
                <h3 className="text-sm font-semibold text-gray-900">{review.clientName}</h3>
                <div className="flex items-center gap-2 mt-1">
                  <StarRow rating={review.rating} />
                  <span className="text-xs text-gray-400">{formatDate(review.createdAt)}</span>
                </div>
              </div>
              {review.category && (
                <span className="text-xs font-mono text-gray-400 bg-gray-50 px-2 py-1 rounded">
                  {review.category}
                </span>
              )}
            </div>
            <p className="text-sm text-gray-600 leading-relaxed">{review.comment}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
