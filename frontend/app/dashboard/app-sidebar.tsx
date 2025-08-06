"use client";

import * as React from "react"

import { SearchForm } from "@/app/dashboard/search-form"
import { VersionSwitcher } from "@/components/version-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"
import { useEffect, useState } from "react";
import Link from 'next/link'
import { NextRequest } from 'next/server';
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { Input } from "@/components/ui/input";

type Project = {
  id: string;
  title: string;
  isActive: boolean;
};

// type AppSidebarProps = React.ComponentProps<typeof Sidebar> & {
//   selectedProject: Project | null;                    // your custom prop
//   setSelectedProject: (project: Project | null) => void;   // your custom setter
// };

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const [projects, setProjects] = useState<Project[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [newProjectName, setNewProjectName] = useState("");

  async function fetchProjects() {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("/api/projects", {
        method: "GET",
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error(`Error ${res.status}`);
      }
      const data = await res.json();
      setProjects(data.projects);
    } catch (err) {
      setError("Failed to load projects");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchProjects();
  }, []);

  function handleProjectClick(clickedId: string) {
    setProjects((prev) =>
      prev
        ? prev.map((project) => ({
          ...project,
          isActive: project.id === clickedId,
        }))
        : null
    );
  }

  async function handleCreateProject(projectName: string) {
    if (!projectName.trim()) {
      setError("Project name cannot be empty");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const res = await fetch("/api/projects/add-project", {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: projectName }),
      });

      if (!res.ok) {
        const errorData = await res.json();
        const message = errorData.error || errorData.message || "Failed to create project";
        setError(message);
        return;
      }

      setNewProjectName("");

      // refetch projects to update the list
      await fetchProjects();
    } catch (err) {
      setError("Network error: Failed to create project");
    } finally {
      setLoading(false);
    }
  }


  return (
    <Sidebar {...props}>
      <SidebarHeader>
        {/* <VersionSwitcher
          versions={data.versions}
          defaultVersion={data.versions[0]}
        /> */}
        <SearchForm />
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup key="projects">
          <SidebarGroupLabel>Projects</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {loading && (
                <SidebarMenuItem>
                  <SidebarMenuButton className="text-gray-600">
                    Loading...
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )}

              {!loading && (!projects || projects.length === 0) && (
                <SidebarMenuItem>
                  <SidebarMenuButton className="text-gray-600">
                    No projects yet... 😴
                  </SidebarMenuButton>
                </SidebarMenuItem>
              )}

              {!loading && projects && projects.length > 0 && projects.map((item) => (
                <SidebarMenuItem key={item.id}>
                  <SidebarMenuButton asChild isActive={item.isActive} onClick={() => handleProjectClick(item.id)}>
                    <Link href={`/dashboard/${item.id}/tasks`}>{item.title}</Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="flex flex-col gap-2 p-4">
        <Input
          type="text"
          placeholder="Enter new project name"
          value={newProjectName}
          onChange={(e) => setNewProjectName(e.target.value)}
          disabled={loading}
        />
        <Button
          variant="outline"
          size="sm"
          onClick={() => handleCreateProject(newProjectName)}
          disabled={loading || newProjectName.trim() === ""}
          className="flex items-center justify-center gap-2"
        >
          Create Project <Plus className="h-4 w-4" />
        </Button>
        {error && <p className="text-red-600 text-sm mt-1">{error}</p>}
      </SidebarFooter>

      <SidebarRail />
    </Sidebar>
  )
}
