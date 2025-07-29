"use client";

import * as React from "react"

import { SearchForm } from "@/app/dashboard/search-form"
import { VersionSwitcher } from "@/components/version-switcher"
import {
  Sidebar,
  SidebarContent,
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

type Project = {
  id: string;
  title: string;
  isActive: boolean;
};


export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {

  const [projects, setProjects] = useState<Project[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchProjects() {
      setLoading(true);
      setError(null);
      try {
        const res = await fetch("/api/projects", {
          method: "GET",
          credentials: "include", // Sends cookies securely (including httpOnly)
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

    fetchProjects();
  }, [])


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
                  <SidebarMenuButton asChild isActive={item.isActive}>
                    <Link href={`dashboard/${item.id}/tasks`}>{item.title}</Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarRail />
    </Sidebar>
  )
}
