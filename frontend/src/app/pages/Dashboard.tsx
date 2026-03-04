import { useState } from "react";
import { RefreshCw, Sparkles, Calendar } from "lucide-react";
import KPICards from "../components/KPICards";
import MyPostsToday from "../components/MyPostsToday";
import NewLikesToday from "../components/NewLikesToday";
import NewRepliesToday from "../components/NewRepliesToday";
import MentionsToday from "../components/MentionsToday";
import TopEngagersToday from "../components/TopEngagersToday";
import * as Tabs from "@radix-ui/react-tabs";

export default function Dashboard() {
  const [isRefreshing, setIsRefreshing] = useState(false);

  const handleRefresh = () => {
    setIsRefreshing(true);
    setTimeout(() => setIsRefreshing(false), 1000);
  };

  const handleGenerateDigest = () => {
    alert("Generating comprehensive digest... (demo)");
  };

  const todayDate = new Date().toLocaleDateString("en-US", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className="p-4 lg:p-8 pb-24 lg:pb-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-2">
          <div>
            <h1 className="text-3xl lg:text-4xl font-bold text-white mb-2">Today's Digest</h1>
            <div className="flex items-center gap-2 text-gray-400">
              <Calendar className="w-4 h-4" />
              <span className="text-sm">{todayDate}</span>
            </div>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={handleRefresh}
              className={`px-4 py-2.5 rounded-[16px] backdrop-blur-xl bg-white/5 border border-white/10 text-white hover:bg-white/10 transition-all flex items-center gap-2 ${
                isRefreshing ? "animate-pulse" : ""
              }`}
            >
              <RefreshCw className={`w-4 h-4 ${isRefreshing ? "animate-spin" : ""}`} />
              <span className="hidden sm:inline">Refresh</span>
            </button>
            <button
              onClick={handleGenerateDigest}
              className="px-6 py-2.5 rounded-[16px] bg-gradient-to-r from-purple-500 to-blue-500 text-white hover:from-purple-600 hover:to-blue-600 transition-all flex items-center gap-2 shadow-lg"
            >
              <Sparkles className="w-4 h-4" />
              <span>Generate Digest</span>
            </button>
          </div>
        </div>
      </div>

      {/* KPI Cards */}
      <KPICards />

      {/* My Posts Today */}
      <div className="mb-6">
        <MyPostsToday />
      </div>

      {/* Desktop: Three column layout for panels */}
      <div className="hidden lg:grid lg:grid-cols-3 gap-6 mb-6">
        <NewLikesToday />
        <NewRepliesToday />
        <MentionsToday />
      </div>

      {/* Mobile: Tabs for panels */}
      <div className="lg:hidden mb-6">
        <Tabs.Root defaultValue="likes" className="w-full">
          <Tabs.List className="flex gap-2 mb-4 backdrop-blur-xl bg-white/5 p-1.5 rounded-[16px] border border-white/10">
            <Tabs.Trigger
              value="likes"
              className="flex-1 px-4 py-2 rounded-[12px] text-sm text-gray-400 data-[state=active]:bg-white/10 data-[state=active]:text-white transition-all"
            >
              Likes
            </Tabs.Trigger>
            <Tabs.Trigger
              value="replies"
              className="flex-1 px-4 py-2 rounded-[12px] text-sm text-gray-400 data-[state=active]:bg-white/10 data-[state=active]:text-white transition-all"
            >
              Replies
            </Tabs.Trigger>
            <Tabs.Trigger
              value="mentions"
              className="flex-1 px-4 py-2 rounded-[12px] text-sm text-gray-400 data-[state=active]:bg-white/10 data-[state=active]:text-white transition-all"
            >
              Mentions
            </Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="likes">
            <NewLikesToday />
          </Tabs.Content>
          <Tabs.Content value="replies">
            <NewRepliesToday />
          </Tabs.Content>
          <Tabs.Content value="mentions">
            <MentionsToday />
          </Tabs.Content>
        </Tabs.Root>
      </div>

      {/* Top Engagers - Full width */}
      <TopEngagersToday />
    </div>
  );
}
