-- 跨表外键，待所有表及候选唯一键创建后执行
-- 按文件编号升序执行；本文件与对应 down 文件由迁移工具记录版本。
BEGIN;

ALTER TABLE users ADD CONSTRAINT fk_users_avatar_asset_id FOREIGN KEY (avatar_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE auth_identities ADD CONSTRAINT fk_auth_identities_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE agreement_acceptances ADD CONSTRAINT fk_agreement_acceptances_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE agreement_acceptances ADD CONSTRAINT fk_agreement_acceptances_content_asset_id FOREIGN KEY (content_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE agreement_acceptances ADD CONSTRAINT fk_agreement_acceptances_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE media_assets ADD CONSTRAINT fk_media_assets_owner_user_id FOREIGN KEY (owner_user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE media_assets ADD CONSTRAINT fk_media_assets_uploaded_by_admin_id FOREIGN KEY (uploaded_by_admin_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE notifications ADD CONSTRAINT fk_notifications_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE coaches ADD CONSTRAINT fk_coaches_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE coach_applications ADD CONSTRAINT fk_coach_applications_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE coach_applications ADD CONSTRAINT fk_coach_applications_identity_check_id FOREIGN KEY (identity_check_id) REFERENCES identity_checks (id) ON DELETE RESTRICT;

ALTER TABLE coach_applications ADD CONSTRAINT fk_coach_applications_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE identity_checks ADD CONSTRAINT fk_identity_checks_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE identity_checks ADD CONSTRAINT fk_identity_checks_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE coach_skills ADD CONSTRAINT fk_coach_skills_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE coach_skills ADD CONSTRAINT fk_coach_skills_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE coach_skills ADD CONSTRAINT fk_coach_skills_rank_id FOREIGN KEY (rank_id) REFERENCES game_ranks (id) ON DELETE RESTRICT;

ALTER TABLE coach_skills ADD CONSTRAINT fk_coach_skills_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE games ADD CONSTRAINT fk_games_icon_asset_id FOREIGN KEY (icon_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE game_regions ADD CONSTRAINT fk_game_regions_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE game_ranks ADD CONSTRAINT fk_game_ranks_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE service_categories ADD CONSTRAINT fk_service_categories_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE coach_services ADD CONSTRAINT fk_coach_services_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE coach_services ADD CONSTRAINT fk_coach_services_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE coach_services ADD CONSTRAINT fk_coach_services_category_id FOREIGN KEY (category_id) REFERENCES service_categories (id) ON DELETE RESTRICT;

ALTER TABLE coach_service_regions ADD CONSTRAINT fk_coach_service_regions_service_id FOREIGN KEY (service_id) REFERENCES coach_services (id) ON DELETE RESTRICT;

ALTER TABLE coach_service_regions ADD CONSTRAINT fk_coach_service_regions_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE coach_service_regions ADD CONSTRAINT fk_coach_service_regions_region_id FOREIGN KEY (region_id) REFERENCES game_regions (id) ON DELETE RESTRICT;

ALTER TABLE coach_availability ADD CONSTRAINT fk_coach_availability_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE booking_holds ADD CONSTRAINT fk_booking_holds_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE booking_holds ADD CONSTRAINT fk_booking_holds_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE order_quotes ADD CONSTRAINT fk_order_quotes_buyer_id FOREIGN KEY (buyer_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE order_quotes ADD CONSTRAINT fk_order_quotes_service_id FOREIGN KEY (service_id) REFERENCES coach_services (id) ON DELETE RESTRICT;

ALTER TABLE order_quotes ADD CONSTRAINT fk_order_quotes_region_id FOREIGN KEY (region_id) REFERENCES game_regions (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_quote_id FOREIGN KEY (quote_id) REFERENCES order_quotes (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_buyer_id FOREIGN KEY (buyer_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_service_id FOREIGN KEY (service_id) REFERENCES coach_services (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_region_id FOREIGN KEY (region_id) REFERENCES game_regions (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_commission_rule_id FOREIGN KEY (commission_rule_id) REFERENCES commission_rules (id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_business_rule_id FOREIGN KEY (business_rule_id) REFERENCES business_rule_versions (id) ON DELETE RESTRICT;

ALTER TABLE order_events ADD CONSTRAINT fk_order_events_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE payments ADD CONSTRAINT fk_payments_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE payment_events ADD CONSTRAINT fk_payment_events_payment_id FOREIGN KEY (payment_id) REFERENCES payments (id) ON DELETE RESTRICT;

ALTER TABLE payment_events ADD CONSTRAINT fk_payment_events_payload_asset_id FOREIGN KEY (payload_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_payment_id FOREIGN KEY (payment_id) REFERENCES payments (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_applicant_user_id FOREIGN KEY (applicant_user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_case_id FOREIGN KEY (case_id) REFERENCES after_sales_cases (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_approval_id FOREIGN KEY (approval_id) REFERENCES admin_approvals (id) ON DELETE RESTRICT;

ALTER TABLE refund_events ADD CONSTRAINT fk_refund_events_refund_id FOREIGN KEY (refund_id) REFERENCES refunds (id) ON DELETE RESTRICT;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_buyer_id FOREIGN KEY (buyer_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE reviews ADD CONSTRAINT fk_reviews_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE review_contents ADD CONSTRAINT fk_review_contents_review_id FOREIGN KEY (review_id) REFERENCES reviews (id) ON DELETE RESTRICT;

ALTER TABLE after_sales_cases ADD CONSTRAINT fk_after_sales_cases_reporter_id FOREIGN KEY (reporter_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE after_sales_cases ADD CONSTRAINT fk_after_sales_cases_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE after_sales_cases ADD CONSTRAINT fk_after_sales_cases_target_user_id FOREIGN KEY (target_user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE after_sales_cases ADD CONSTRAINT fk_after_sales_cases_review_id FOREIGN KEY (review_id) REFERENCES reviews (id) ON DELETE RESTRICT;

ALTER TABLE after_sales_cases ADD CONSTRAINT fk_after_sales_cases_assignee_id FOREIGN KEY (assignee_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE case_events ADD CONSTRAINT fk_case_events_case_id FOREIGN KEY (case_id) REFERENCES after_sales_cases (id) ON DELETE RESTRICT;

ALTER TABLE settlement_accounts ADD CONSTRAINT fk_settlement_accounts_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_category_id FOREIGN KEY (category_id) REFERENCES service_categories (id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_created_by FOREIGN KEY (created_by) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_approval_id FOREIGN KEY (approval_id) REFERENCES admin_approvals (id) ON DELETE RESTRICT;

ALTER TABLE settlement_items ADD CONSTRAINT fk_settlement_items_order_id FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE RESTRICT;

ALTER TABLE settlement_items ADD CONSTRAINT fk_settlement_items_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE settlement_items ADD CONSTRAINT fk_settlement_items_account_id FOREIGN KEY (account_id) REFERENCES settlement_accounts (id) ON DELETE RESTRICT;

ALTER TABLE settlement_adjustments ADD CONSTRAINT fk_settlement_adjustments_settlement_item_id FOREIGN KEY (settlement_item_id) REFERENCES settlement_items (id) ON DELETE RESTRICT;

ALTER TABLE settlement_adjustments ADD CONSTRAINT fk_settlement_adjustments_refund_id FOREIGN KEY (refund_id) REFERENCES refunds (id) ON DELETE RESTRICT;

ALTER TABLE settlement_adjustments ADD CONSTRAINT fk_settlement_adjustments_approval_id FOREIGN KEY (approval_id) REFERENCES admin_approvals (id) ON DELETE RESTRICT;

ALTER TABLE settlement_batches ADD CONSTRAINT fk_settlement_batches_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE settlement_batches ADD CONSTRAINT fk_settlement_batches_account_id FOREIGN KEY (account_id) REFERENCES settlement_accounts (id) ON DELETE RESTRICT;

ALTER TABLE settlement_batch_items ADD CONSTRAINT fk_settlement_batch_items_batch_id FOREIGN KEY (batch_id) REFERENCES settlement_batches (id) ON DELETE RESTRICT;

ALTER TABLE settlement_batch_items ADD CONSTRAINT fk_settlement_batch_items_settlement_item_id FOREIGN KEY (settlement_item_id) REFERENCES settlement_items (id) ON DELETE RESTRICT;

ALTER TABLE settlement_batch_items ADD CONSTRAINT fk_settlement_batch_items_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE reconciliation_runs ADD CONSTRAINT fk_reconciliation_runs_bill_asset_id FOREIGN KEY (bill_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE reconciliation_lines ADD CONSTRAINT fk_reconciliation_lines_run_id FOREIGN KEY (run_id) REFERENCES reconciliation_runs (id) ON DELETE RESTRICT;

ALTER TABLE reconciliation_diffs ADD CONSTRAINT fk_reconciliation_diffs_run_id FOREIGN KEY (run_id) REFERENCES reconciliation_runs (id) ON DELETE RESTRICT;

ALTER TABLE reconciliation_diffs ADD CONSTRAINT fk_reconciliation_diffs_line_id FOREIGN KEY (line_id) REFERENCES reconciliation_lines (id) ON DELETE RESTRICT;

ALTER TABLE reconciliation_diffs ADD CONSTRAINT fk_reconciliation_diffs_assignee_id FOREIGN KEY (assignee_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE activities ADD CONSTRAINT fk_activities_cover_asset_id FOREIGN KEY (cover_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE activities ADD CONSTRAINT fk_activities_created_by FOREIGN KEY (created_by) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE banners ADD CONSTRAINT fk_banners_image_asset_id FOREIGN KEY (image_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE coach_recommendations ADD CONSTRAINT fk_coach_recommendations_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE coach_recommendations ADD CONSTRAINT fk_coach_recommendations_game_id FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE RESTRICT;

ALTER TABLE content_audits ADD CONSTRAINT fk_content_audits_submitted_by_user_id FOREIGN KEY (submitted_by_user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE content_audits ADD CONSTRAINT fk_content_audits_submitted_by_admin_id FOREIGN KEY (submitted_by_admin_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE content_audits ADD CONSTRAINT fk_content_audits_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE user_penalties ADD CONSTRAINT fk_user_penalties_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE user_penalties ADD CONSTRAINT fk_user_penalties_coach_id FOREIGN KEY (coach_id) REFERENCES coaches (id) ON DELETE RESTRICT;

ALTER TABLE user_penalties ADD CONSTRAINT fk_user_penalties_service_id FOREIGN KEY (service_id) REFERENCES coach_services (id) ON DELETE RESTRICT;

ALTER TABLE user_penalties ADD CONSTRAINT fk_user_penalties_case_id FOREIGN KEY (case_id) REFERENCES after_sales_cases (id) ON DELETE RESTRICT;

ALTER TABLE user_penalties ADD CONSTRAINT fk_user_penalties_admin_id FOREIGN KEY (admin_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE admin_user_roles ADD CONSTRAINT fk_admin_user_roles_admin_id FOREIGN KEY (admin_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE admin_user_roles ADD CONSTRAINT fk_admin_user_roles_role_id FOREIGN KEY (role_id) REFERENCES admin_roles (id) ON DELETE RESTRICT;

ALTER TABLE admin_role_permissions ADD CONSTRAINT fk_admin_role_permissions_role_id FOREIGN KEY (role_id) REFERENCES admin_roles (id) ON DELETE RESTRICT;

ALTER TABLE admin_role_permissions ADD CONSTRAINT fk_admin_role_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES admin_permissions (id) ON DELETE RESTRICT;

ALTER TABLE admin_approvals ADD CONSTRAINT fk_admin_approvals_requester_id FOREIGN KEY (requester_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE admin_approvals ADD CONSTRAINT fk_admin_approvals_reviewer_id FOREIGN KEY (reviewer_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE admin_audit_logs ADD CONSTRAINT fk_admin_audit_logs_admin_id FOREIGN KEY (admin_id) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE admin_audit_logs ADD CONSTRAINT fk_admin_audit_logs_payload_asset_id FOREIGN KEY (payload_asset_id) REFERENCES media_assets (id) ON DELETE RESTRICT;

ALTER TABLE business_rule_versions ADD CONSTRAINT fk_business_rule_versions_created_by FOREIGN KEY (created_by) REFERENCES admin_users (id) ON DELETE RESTRICT;

ALTER TABLE business_rule_versions ADD CONSTRAINT fk_business_rule_versions_approval_id FOREIGN KEY (approval_id) REFERENCES admin_approvals (id) ON DELETE RESTRICT;

ALTER TABLE user_activity_events ADD CONSTRAINT fk_user_activity_events_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT;

ALTER TABLE coach_skills ADD CONSTRAINT fk_coach_skills_rank_id_game_id FOREIGN KEY (rank_id, game_id) REFERENCES game_ranks (id, game_id) ON DELETE RESTRICT;

ALTER TABLE coach_services ADD CONSTRAINT fk_coach_services_category_id_game_id FOREIGN KEY (category_id, game_id) REFERENCES service_categories (id, game_id) ON DELETE RESTRICT;

ALTER TABLE coach_service_regions ADD CONSTRAINT fk_coach_service_regions_service_id_game_id FOREIGN KEY (service_id, game_id) REFERENCES coach_services (id, game_id) ON DELETE RESTRICT;

ALTER TABLE coach_service_regions ADD CONSTRAINT fk_coach_service_regions_region_id_game_id FOREIGN KEY (region_id, game_id) REFERENCES game_regions (id, game_id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_region_id_game_id FOREIGN KEY (region_id, game_id) REFERENCES game_regions (id, game_id) ON DELETE RESTRICT;

ALTER TABLE orders ADD CONSTRAINT fk_orders_service_id_game_id_coach_id FOREIGN KEY (service_id, game_id, coach_id) REFERENCES coach_services (id, game_id, coach_id) ON DELETE RESTRICT;

ALTER TABLE booking_holds ADD CONSTRAINT fk_booking_holds_order_id_coach_id FOREIGN KEY (order_id, coach_id) REFERENCES orders (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE refunds ADD CONSTRAINT fk_refunds_payment_id_order_id FOREIGN KEY (payment_id, order_id) REFERENCES payments (id, order_id) ON DELETE RESTRICT;

ALTER TABLE settlement_items ADD CONSTRAINT fk_settlement_items_order_id_coach_id FOREIGN KEY (order_id, coach_id) REFERENCES orders (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_items ADD CONSTRAINT fk_settlement_items_account_id_coach_id FOREIGN KEY (account_id, coach_id) REFERENCES settlement_accounts (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_batches ADD CONSTRAINT fk_settlement_batches_account_id_coach_id FOREIGN KEY (account_id, coach_id) REFERENCES settlement_accounts (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_batch_items ADD CONSTRAINT fk_settlement_batch_items_batch_id_coach_id FOREIGN KEY (batch_id, coach_id) REFERENCES settlement_batches (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_batch_items ADD CONSTRAINT fk_settlement_batch_items_settlement_item_id_coach_id FOREIGN KEY (settlement_item_id, coach_id) REFERENCES settlement_items (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE commission_rules ADD CONSTRAINT fk_commission_rules_category_id_game_id FOREIGN KEY (category_id, game_id) REFERENCES service_categories (id, game_id) ON DELETE RESTRICT;

ALTER TABLE settlement_supplements ADD CONSTRAINT fk_supplement_adjustment
 FOREIGN KEY (adjustment_id, settlement_item_id) REFERENCES settlement_adjustments (id, settlement_item_id) ON DELETE RESTRICT;

ALTER TABLE settlement_supplements ADD CONSTRAINT fk_supplement_item_coach
 FOREIGN KEY (settlement_item_id, coach_id) REFERENCES settlement_items (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_supplements ADD CONSTRAINT fk_supplement_account_coach
 FOREIGN KEY (account_id, coach_id) REFERENCES settlement_accounts (id, coach_id) ON DELETE RESTRICT;

ALTER TABLE settlement_supplements ADD CONSTRAINT fk_supplement_approval
 FOREIGN KEY (approval_id) REFERENCES admin_approvals (id) ON DELETE RESTRICT;

COMMIT;
